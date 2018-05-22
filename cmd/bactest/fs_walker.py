from typing import List, Dict
import os
from post_actions import PostActions, POST_ACTION_SUFFIX
from pre_actions import PreActions, PRE_ACTION_SUFFIX
from tests import Tests, TESTS_SUFFIX
from builder import BactestBuilder
from module import Module
from displayer import Framework

class TestFSWalker:
    def __init__(self, tests_location: str):
        self.__tests_location = os.path.abspath(tests_location)
        self.displayer = Framework(module_name=__class__.__name__)

    @property
    def get_tests_location(self):
        return self.__tests_location

    def read_tests_dirs_names(self) -> Dict:
        self.displayer.info("Looking for tests dirs in: " + self.__tests_location)
        current_path = os.getcwd()
        abs_tests_dirs = {}
        dirs = os.listdir(self.__tests_location)
        for d in dirs:
            abs_tests_dirs[d] = os.path.join(self.__tests_location, d)
        return abs_tests_dirs

    def read_modules_from_dirs(self, dirs: Dict) -> Dict:
        files_dir_map = {}
        for k, v in dirs.items():
            module_files = []
            for f in os.listdir(v):
                if f.endswith("py") and not f.startswith("__"):
                    module_files.append(os.path.join(v, f))
            files_dir_map[k] = module_files
        return files_dir_map

    def find_pre_actions(self, files_dir_map: Dict)->PreActions:
        pre_actions = PreActions()
        for k in files_dir_map.keys():
            if k.startswith(PRE_ACTION_SUFFIX):
                if len(files_dir_map[k]) != 0:
                    dir_name = os.path.dirname(files_dir_map[k][0])
                    pre_action = Module(k, dir_name , files_dir_map[k])
                    pre_actions.modules.append(pre_action)
                else:
                    self.displayer.warning("Skipping module -->{}<-- does not have any tests".format(k))
        return pre_actions

    def find_post_actions(self, files_dir_map: Dict)->PostActions:
        post_actions = PostActions()
        for k in files_dir_map.keys():
            if k.startswith(POST_ACTION_SUFFIX):
                if len(files_dir_map[k]) != 0:
                    dir_name = os.path.dirname(files_dir_map[k][0])
                    post_action = Module(k, dir_name ,files_dir_map[k])
                    post_actions.modules.append(post_action)
                else:
                    self.displayer.warning("Skipping module -->{}<-- does not have any tests".format(k))
        return post_actions


    def find_tests(self, files_dir_map: Dict)->Tests:
        tests = Tests() 
        for k in files_dir_map.keys():
            if k.startswith(TESTS_SUFFIX):
                if len(files_dir_map[k]) != 0:
                    dir_name = os.path.dirname(files_dir_map[k][0])
                    test = Module(k, dir_name, files_dir_map[k])
                    tests.modules.append(test)
                else:
                     self.displayer.warning("Skipping module -->{}<-- does not have any tests".format(k))
        return tests

    def read_modules(self) -> BactestBuilder:
        tests_dirs_names = self.read_tests_dirs_names()
        tests_dirs = self.read_modules_from_dirs(tests_dirs_names)
        tests = self.find_tests(tests_dirs)
        pre_actions = self.find_pre_actions(tests_dirs)
        post_actions = self.find_post_actions(tests_dirs)
        tests_with_actions = BactestBuilder(pre_actions, post_actions, tests)
        return tests_with_actions


