#!/usr/bin/env python3
"""Atajo de S2-02: equivale a `python3 -I oraculo/compara.py s2 [opciones]`."""
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import compara  # noqa: E402

if __name__ == '__main__':
    sys.exit(compara.main(['s2'] + sys.argv[1:]))
