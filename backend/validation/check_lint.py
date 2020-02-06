"""
Script for checking lint checking
"""
from yamllint.config import YamlLintConfig
from yamllint import linter
import sys

requested_file=str(sys.argv[1])

# conf object must be created based on .yamllint rules
conf = YamlLintConfig(file='.yamllint')
f = open(requested_file)
gen = linter.run(f, conf,'.yamllint')
errors = list(gen)
if errors:
    print(errors)
else:
    print("Success")