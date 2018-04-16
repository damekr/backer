import datetime
import argparse
import os
import importlib

pre_action_suffix = "pre_"
post_action_suffix = "post_"
tests_suffix = "test_"

class Bactest:
    def __init__(self, tests_location: str):
        self.test_location = tests_location


class BactestResult:
    def __init__(self):
        self.result = bool
        self.className = dict
        self.testName = str

class BactestSummary:
    def __init__(self):
        self.startTime = datetime.datetime.now()
        self.endTime = str
        self.results = []

class TestWalker:
    def __init__(self, tests_location: str):
        self.test_location = tests_location

    def load_classes(self):
        tests = {}




def parse_args():
    parser = argparse.ArgumentParser(description='Bactest tool for making integration tests')
    parser.add_argument('-d', '--directory', action='store', dest='directory', help='directory of tests', required=True)
    args = parser.parse_args()
    return args



if __name__ == "__main__":
    args = parse_args()
    print(args.directory)
    test_walker = TestWalker(args.directory)
    test_walker.load_classes()

