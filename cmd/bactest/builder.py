import importlib.machinery
from pre_actions import PreActions
from post_actions import PostActions
from displayer import Framework
from tests import Tests
from typing import List
import os


class BactestBuilder:
    def __init__(self, preActions: PreActions, postActions: PostActions, tests: Tests):
        self.preModules = preActions
        self.postModules = postActions
        self.tests = tests
        self.displayer = Framework(module_name=__class__.__name__)

    def import_all(self):
        self.__import_pre_actions()
        self.__import_tests()
        self.__import_post_actions()

    def __import_tests(self):
        for v in self.tests.modules:
            for f in v.files:
                self.displayer.debug("NAME: " + v.name)
                filename = os.path.basename(f)
                self.displayer.debug("Filename:  " + filename)
                self.displayer.debug("FILE: " + f)
                module = importlib.machinery.SourceFileLoader(filename, f).load_module()
                v.executables.append(module)

    def __import_pre_actions(self):
        for v in self.preModules.modules:
            for f in v.files:
                self.displayer.debug("PRE_PATH: " + f)
                filename = os.path.basename(f)
                module = importlib.machinery.SourceFileLoader(filename, f).load_module()
                v.executables.append(module)

    def __import_post_actions(self):
        for v in self.postModules.modules:
            for f in v.files:
                self.displayer.debug("POST_PATH: " + f)
                filename = os.path.basename(f)
                module = importlib.machinery.SourceFileLoader(filename, f).load_module()
                v.executables.append(module)

    def get_tests_files_paths(self) -> List:
        names = []
        for n in self.tests.modules:
            names.append(n.paths)
        return names

    def get_pre_actions_files_paths(self) -> List:
        names = []
        for n in self.preModules.modules:
            names.append(n.paths)
        return names

    def get_post_actions_files_paths(self) -> List:
        names = []
        for n in self.postModules.modules:
            names.append(n.paths)
        return names


