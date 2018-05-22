import argparse
from fs_walker import TestFSWalker
from displayer import Framework
from executor import TestExecutor

class Bactest:
    def __init__(self, tests_location: str):
        self.test_location = tests_location


def parse_args():
    parser = argparse.ArgumentParser(description='Bactest tool for making integration tests')
    parser.add_argument('-d', '--directory', action='store', dest='directory', help='directory of tests', required=True)
    args = parser.parse_args()
    return args


if __name__ == "__main__":
    args = parse_args()
    displayer = Framework(module_name=__name__)
    test_walker = TestFSWalker(args.directory)
    tests_builder = test_walker.read_modules()
    tests_builder.import_all()
    test_executor = TestExecutor()
    test_executor.run(tests_builder)
