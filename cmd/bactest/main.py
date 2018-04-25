import datetime
import argparse
import os
from typing import Dict, List
from importlib.abc import Loader
import re


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

class Module:
    def __init__(self):
        self.path = ""
        self.name = ""



class PreActions:
    def __init__(self):
        self.modules = List[Module]

class PostActions:
    def __init__(self):
        self.modules = List[Module]

class Tests:
    def __init__(self, preActions: PreActions, postActions: PostActions):
        self.preModules = preActions
        self.postModules = postActions
        self.modules = List[Module]



class TestFSWalker(Loader):
    def __init__(self, tests_location: str):
        self.__tests_location = tests_location

    @property
    def getTestsLocation(self):
        return self.__tests_location

    def read_tests_dirs_names(self) -> Dict:
        print("Looking for tests dirs in: ", self.__tests_location)
        current_path = os.getcwd()
        abs_tests_dirs = {}
        print("Current path: ", current_path)
        dirs = os.listdir(self.__tests_location)
        for d in dirs:
            abs_tests_dirs[d] = os.path.join(current_path, d)
        return abs_tests_dirs

    def read_tests_packages(self, abs_tests_dirs: Dict) -> Dict:
        modules_dict = {}
        test = re.compile("test\.py$", re.IGNORECASE)
        for name, path in abs_tests_dirs.items():
            print("Importing module: ", name)
            print("From path: ", path)
            modules_dict[name] = self.load_module(path)
        return modules_dict





def parse_args():
    parser = argparse.ArgumentParser(description='Bactest tool for making integration tests')
    parser.add_argument('-d', '--directory', action='store', dest='directory', help='directory of tests', required=True)
    args = parser.parse_args()
    return args



if __name__ == "__main__":
    args = parse_args()
    print(args.directory)
    test_walker = TestFSWalker(args.directory)
    tests_dirs = test_walker.read_tests_dirs_names()
    for k, v in test_walker.read_tests_packages(tests_dirs).items():
        print(type(k))
    print(globals())
