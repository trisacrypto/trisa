package trust

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"software.sslmate.com/src/go-pkcs12"
)

// Compression
const (
	CompressionGZIP = ".gz"
	CompressionZIP  = ".zip"
	CompressionNone = ".pem"
	CompressionAuto = "auto"
)

var validFormats = map[string]struct{}{
	CompressionGZIP: {},
	CompressionZIP:  {},
	CompressionNone: {},
	CompressionAuto: {},
}

// Serializer maintains options for compression, encoding, and pkcs12 encryption when
// serializing and deserializing Provider and ProviderPool objects to and from disk.
type Serializer struct {
	// During extraction, if marked as private the serializer will expect PKCS12
	// decryption from the incoming data. During compression, if marked not private,
	// then the serializer will only serialize the public Provider. This has no effect
	// on ProviderPool serialization which can only be serialized/deserialized with the
	// public certificates.
	Private bool

	// The password for encryption or decryption; required if Private is set to true,
	// if empty then the pkcs12.DefaultPassword, "changeit" is used.
	Password string

	// Format is one of "gz", "zip", "pem" (for no compression) and auto. If auto, then
	// during extraction the format is detected from the file extension, otherwise an
	// error is returned if directly from bytes. During compression if auto, then the
	// format is "gz" for Providers and "zip" for ProviderPools.
	Format string

	// Internal helper fields
	path     string
	multiple bool
}

// NewSerializer creates a serializer with default options ready to either extract or
// compress a Provider or a ProviderPool. The serializer takes additional positional
// arguments of password and format respectively.
func NewSerializer(private bool, opts ...string) (_ *Serializer, err error) {
	serializer := &Serializer{Private: private}

	if len(opts) > 0 {
		serializer.Password = opts[0]
	}
	if serializer.Private && serializer.Password == "" {
		serializer.Password = pkcs12.DefaultPassword
	}

	if len(opts) > 1 {
		serializer.Format = opts[1]
		if _, err = serializer.getFormat(); err != nil {
			return nil, err
		}
	} else {
		serializer.Format = CompressionAuto
	}

	return serializer, nil
}

// Extract a Provider object from the given data, decompressing according to the format
// and decoding or decrypting the data as needed. Returns an error if the format is not
// correctly set or available.
func (s *Serializer) Extract(data []byte) (p *Provider, err error) {
	b := bytes.NewReader(data)
	return s.Read(b)
}

// ExtractPool from the given data, decompressing according to the format and decoding
// or decrypting the data as needed. Returns an error if the format is not correctly set
// or available on the serializer object.
func (s *Serializer) ExtractPool(data []byte) (pool ProviderPool, err error) {
	b := bytes.NewReader(data)
	return s.ReadPool(b)
}

// Read and extract a Provider from the reader object. If the Format is CompressionZip,
// this method expects an io.ReadCloser, returned from the zip.File.Open method. Archive
// handling must happen before this method is called.
func (s *Serializer) Read(r io.Reader) (_ *Provider, err error) {
	// Set internal fields for downstream computation
	s.multiple = false
	var mode string
	if mode, err = s.getFormat(); err != nil {
		return nil, err
	}

	// Load the data into memory
	var data []byte

	switch mode {
	case CompressionGZIP:
		if r, err = gzip.NewReader(r); err != nil {
			return nil, err
		}
		fallthrough
	case CompressionNone, CompressionZIP:
		if data, err = ioutil.ReadAll(r); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unhandled format %q", mode)
	}

	if s.Private {
		return Decrypt(data, s.Password)
	}

	return New(data)
}

// ReadPool and extract a ProviderPool from the reader objects. If no readers are
// provided then an empty pool is returned. If the Format is CompressionZip, this method
// expects that the archive has been opened and is operating on all internal readers.
func (s *Serializer) ReadPool(readers ...io.Reader) (pool ProviderPool, err error) {
	pool = make(ProviderPool)
	for _, r := range readers {
		var p *Provider
		if p, err = s.Read(r); err != nil {
			return nil, err
		}
		pool.Add(p)
	}
	return pool, nil
}

// ReadFile and extract a Provider from it. This function expects only a single provider
// so it handles Zip archives specially, requiring that the Zip archive only has a
// single file in it, otherwise an error is returned.
func (s *Serializer) ReadFile(path string) (p *Provider, err error) {
	// Set internal fields for downstream computation
	s.path = path

	// Open ZIP archive for reading if the extension is .zip
	var mode string
	if mode, err = s.getFormat(); err != nil {
		return nil, err
	}

	if mode == CompressionZIP {
		var archive *zip.ReadCloser
		if archive, err = zip.OpenReader(path); err != nil {
			return nil, err
		}
		defer archive.Close()

		switch len(archive.File) {
		case 0:
			return nil, ErrZipEmpty
		case 1:
			var rc io.ReadCloser
			if rc, err = archive.File[0].Open(); err != nil {
				return nil, err
			}
			defer rc.Close()
			return s.Read(rc)
		default:
			return nil, ErrZipTooMany
		}
	}

	// Handle all other file types
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return nil, err
	}

	return s.Read(f)
}

// ReadPoolFile and extract a ProviderPool from it. This method primarily expects a Zip
// archive with multiple public provider files contained.
// TODO: handle directories
func (s *Serializer) ReadPoolFile(path string) (p ProviderPool, err error) {
	// Set internal fields for downstream computation
	s.path = path

	// Open ZIP archive for reading if the extension is .zip
	var mode string
	if mode, err = s.getFormat(); err != nil {
		return nil, err
	}

	if mode == CompressionZIP {
		var archive *zip.ReadCloser
		if archive, err = zip.OpenReader(path); err != nil {
			return nil, err
		}
		defer archive.Close()

		readers := make([]io.Reader, 0, len(archive.File))
		for _, f := range archive.File {
			var rc io.ReadCloser
			if rc, err = f.Open(); err != nil {
				return nil, err
			}
			defer rc.Close()
			readers = append(readers, rc)
		}
		return s.ReadPool(readers...)
	}

	// Handle all other file types
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return nil, err
	}

	return s.ReadPool(f)
}

// Compress a Provider into the given format with the specified encryption or encoding.
func (s *Serializer) Compress(p *Provider) (data []byte, err error) {
	// Set internal fields for downstream computation
	s.multiple = false

	var b bytes.Buffer
	if err = s.Write(p, &b); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// CompressPool into the given format with the specified encryption or encoding.
func (s *Serializer) CompressPool(pool ProviderPool) (data []byte, err error) {
	// Set internal fields for downstream computation
	s.multiple = len(pool) > 1

	var b bytes.Buffer
	if err = s.WritePool(pool, &b); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Write the provider encoded to the writer object.
func (s *Serializer) Write(p *Provider, w io.Writer) (err error) {
	// Set internal fields for downstream computation
	s.multiple = false
	var mode string
	if mode, err = s.getFormat(); err != nil {
		return err
	}

	// Serialize the data
	var data []byte
	if s.Private {
		if data, err = p.Encrypt(s.Password); err != nil {
			return err
		}
	} else {
		if data, err = p.Encode(); err != nil {
			return err
		}
	}

	// Compress and write the data to disk
	switch mode {
	case CompressionGZIP:
		archive := gzip.NewWriter(w)
		if _, err = archive.Write(data); err != nil {
			return err
		}
		return archive.Close()
	case CompressionZIP:
		archive := zip.NewWriter(w)
		name := p.String() + ".pem"
		name = strings.ReplaceAll(strings.ToLower(name), " ", "_")

		var rc io.Writer
		if rc, err = archive.Create(name); err != nil {
			return err
		}
		if _, err = rc.Write(data); err != nil {
			return err
		}
		return archive.Close()
	case CompressionNone:
		_, err = w.Write(data)
		return err
	default:
		return fmt.Errorf("unhandled format %q", mode)
	}
}

// WritePool encoded to the writer object.
func (s *Serializer) WritePool(pool ProviderPool, w io.Writer) (err error) {
	// Set internal fields for downstream computation
	s.multiple = len(pool) > 1

	// getFormat checks to ensure that we're not writing multiple providers to a non-Zip file
	var mode string
	if mode, err = s.getFormat(); err != nil {
		return err
	}

	// Compress and write the data to disk
	switch mode {
	case CompressionGZIP:
		archive := gzip.NewWriter(w)
		for _, p := range pool {
			// Only write public pools to disk
			p = p.Public()
			var data []byte
			if data, err = p.Encode(); err != nil {
				return err
			}
			if _, err = archive.Write(data); err != nil {
				return err
			}
		}
		return archive.Close()
	case CompressionZIP:
		archive := zip.NewWriter(w)
		for _, p := range pool {
			// Only write public pools to disk
			p = p.Public()
			var data []byte
			if data, err = p.Encode(); err != nil {
				return err
			}

			name := p.String() + ".pem"
			name = strings.ReplaceAll(strings.ToLower(name), " ", "_")

			var rc io.Writer
			if rc, err = archive.Create(name); err != nil {
				return err
			}
			if _, err = rc.Write(data); err != nil {
				return err
			}
		}
		return archive.Close()
	case CompressionNone:
		for _, p := range pool {
			// Only write public pools to disk
			p = p.Public()
			var data []byte
			if data, err = p.Encode(); err != nil {
				return err
			}
			if _, err = w.Write(data); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unhandled format %q", mode)
	}

	return nil
}

// WriteFile with the encoded provider object.
func (s *Serializer) WriteFile(p *Provider, path string) (err error) {
	// Set internal fields for downstream computation
	s.path = path
	s.multiple = false

	var f *os.File
	if f, err = os.Create(path); err != nil {
		return err
	}

	return s.Write(p, f)
}

// WritePoolFile with the encoded provider pool object.
func (s *Serializer) WritePoolFile(pool ProviderPool, path string) (err error) {
	// Set internal fields for downstream computation
	s.path = path
	s.multiple = len(pool) > 1

	var f *os.File
	if f, err = os.Create(path); err != nil {
		return err
	}

	return s.WritePool(pool, f)
}

func (s *Serializer) getFormat() (string, error) {
	if s.Format == "" {
		s.Format = CompressionAuto
	} else {
		s.Format = strings.ToLower(s.Format)
	}

	if _, ok := validFormats[s.Format]; !ok {
		return "", fmt.Errorf("%q not a valid serializer format", s.Format)
	}

	var ext string
	if s.path != "" {
		ext = filepath.Ext(s.path)
	}

	if s.Format == CompressionAuto {
		if ext != "" {
			if _, ok := validFormats[ext]; !ok {
				return "", fmt.Errorf("%q is not a valid serializer format", s.Format)
			}
			return ext, nil
		}

		if s.multiple {
			return CompressionZIP, nil
		}
		return CompressionGZIP, nil
	}

	if ext != "" && ext != s.Format {
		return "", fmt.Errorf("cannot write %s format to %s", s.Format, s.path)
	}

	if s.multiple && s.Format != CompressionZIP {
		return "", fmt.Errorf("cannot write multiple providers to %s", s.Format)
	}

	return s.Format, nil
}
