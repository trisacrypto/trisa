#!/usr/bin/env python3

import os
import glob
import json

from datetime import datetime

CODE_FILE = "dti.gen.go"
DATA_GLOB = "testdata/DTI_Data_*.json"
HEADERS = [
    "// Code generated by generate.py. DO NOT EDIT.",
    "// source: %s",
    "",
    "package dti",
    "",
    "var identifiers = []DTI{",
    "",
]


def load_dti_data(path=None):
    if path is None:
        path = find_dti_data()

    with open(path, 'r') as f:
        data = json.load(f)
        return data["records"]


def find_dti_data(pattern=DATA_GLOB):
    latest = None
    fname = None

    for path in glob.glob(pattern):
        basename, _ = os.path.splitext(os.path.basename(path))
        parts = basename.split("_")
        ts = datetime.strptime(parts[-1], "%Y%m%d")

        if latest is None or ts > latest:
            latest = ts
            fname = path

    return fname


def main():
    path = find_dti_data()
    records = load_dti_data(path)

    with open(CODE_FILE, "w") as f:
        f.write("\n".join(HEADERS) % path)

        for record in records:
            header = record.get("Header")
            info = record.get("Informative")
            shortName = "nil"

            names = info.get("ShortNames")
            if names is not None and len(names) > 0:
                shortName = ", ".join([f'"{n["ShortName"]}"' for n in names])
                shortName = "[]string{" + shortName + "}"

            f.write("\t{\n")
            f.write(f"\t\tIdentifier: \"{header.get('DTI')}\",\n")
            f.write(f"\t\tType: {header.get('DTIType')},\n")
            f.write(f"\t\tVersion: \"{header.get('templateVersion')}\",\n")
            f.write(f"\t\tLongName: \"{info.get('LongName')}\",\n")
            f.write(f"\t\tShortNames: {shortName},\n")
            f.write("\t},\n")

        f.write("}\n")


if __name__ == "__main__":
    main()
