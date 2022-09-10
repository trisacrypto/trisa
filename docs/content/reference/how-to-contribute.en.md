---
title: How to Contribute
date: 2022-06-27T12:48:02-04:00
lastmod: 2022-08-10T16:33:48-04:00
description: "How to Contribute"
weight: 40
---

TRISA is an open source project and welcomes contributions!

The TRISA Protocol implementation is hosted on GitHub at [https://github.com/trisacrypto/trisa](https://github.com/trisacrypto/trisa).

The codebase for the Global Directory Service is hosted on GitHub at [https://github.com/trisacrypto/directory](https://github.com/trisacrypto/trisa).

The codebase for the TRISA TestNet is hosted on GitHub at  [https://github.com/trisacrypto/testnet](https://github.com/trisacrypto/testnet).


## Setting up GitHub for Contributing

1. Set up git signing.

    Create GPG (GNU Privacy Guard) keys to enable git signing to enable verified commits. Install [the GPG command line tools](https://www.gnupg.org/download/).

    Generate a GPG key pair, and the GPG key must use RSA with a key size of 4096 bits. (At the prompt, specify the kind of key you want, or press Enter to accept the default RSA and RSA.)
    ```
    $ gpg --full-generate-key

    ```
    View your generated key as follows:
    ```
    $ gpg --list-secret-keys --keyid-format=long
    ```
    Get your full key as follows:
    ```
    $ gpg --armor --export <your key>
    ```
    Copy your GPG key, which looks like the following:

    ```
    -----BEGIN PGP PUBLIC KEY BLOCK-----
    ....................................
    -----END PGP PUBLIC KEY BLOCK-----
    ```
    Add the GPG key in Github by accessing Accounts Settings. And in the sidebar, under **Access**, click **SSH and GPG Keys**. Then click **New GPG key**, and paste your copied GPG key here in the **Key** field.

    After your GPG key has been added, copy the **Key ID** and set up automatic commit signing.
    ```
    $ git config --global commit.gpgsign true
    $ git config --global user.signingkey <Key ID>
    ```

2. Clone the target repository locally and create a feature branch to hold changes.
    ```
    $ git clone https://github.com/trisacrypto/trisa.git
    $ cd trisa
    $ git checkout -b feature-branch-name
    ```
3. After making changes and pushing commits to your feature branch in the target remote repository, open a Pull Request (PR) and request a review from a maintainer.

## Contributing to the Documentation
1. To contribute to the documentation, you need to [install hugo](https://gohugo.io/getting-started/installing/).
    On a Mac, we recommend using Homebrew as follows:
    ```
    $ brew install hugo
    ```
    You can check the version of Hugo as follows:
    ```
    $ hugo version
    ```
    Ensure the installed version of Hugo matches the deployed Hugo version, which can be found [here](https://github.com/trisacrypto/trisa/blob/main/.github/workflows/publish-docs.yml).


2. After making changes to the documentation, review changes by running a local server while in the [`docs`](https://github.com/trisacrypto/trisa/tree/main/docs) subdirectory. The `-D` renders content marked as draft.

    ```
    $ hugo serve -D
    ```
    Navigate to http://localhost:1313/ to view the rendered website with your changes. The website will be re-rendered every time you make changes by refreshing the page.

## Common Issues While Contributing

1. **Failing to sign commits**

    While trying to commit changes, you receive this error below:
    ```
    gpg failed to sign the data fatal: failed to write commit object
    ```
    **Likely Solution:**
    There may have been another gpg key associated with your GitHub account, and the system still recognizes that previous key. Therefore, you may need to stop the `gpg-agent` by doing the following:
    ```
    $ gpgconf --kill gpg-agent
    ```

2. **Error building site due to new commits with shortcode**

    While trying to run `hugo serve -D` after making changes to the TRISA documentation, you receive this error below:
    ```
    Error: Error building site: ".../trisa/docs/content/...": failed to extract shortcode: template for shortcode "shortcode_name" not found
    ```
     **Likely Solution:**
    You may need to register submodules and submodules within. This can be done with the following:
    ```
    git submodule update --init --recursive
    ```