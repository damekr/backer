from builder import BactestBuilder
from displayer import Framework
from pre_actions import PreActions
from post_actions import PostActions
from tests import Tests

class TestExecutor:
    def __init__(self, tests_parallel=False):
        self.tests_parallel = tests_parallel
        self.displayer = Framework(module_name=__class__.__name__)

    def run(self, tests_builder: BactestBuilder):
        self.displayer.info("Starting tests...")
        self.__run_pre_actions(tests_builder.preModules)
        self.__run_tests(tests_builder.tests)
        self.__run_post_actions(tests_builder.postModules)

    def __run_tests(self, main_tests: Tests):
        self.displayer.info("Found {} tests modules".format(len(main_tests)))
        for k in main_tests.modules:
            print(k.run())
         
    
    def __run_test_module_setup(self):
        pass

    def __run_pre_actions(self, pre_actions: PreActions):
        self.displayer.info("Found {} pre actions modules".format(len(pre_actions)))

    def __run_post_actions(self, post_actions: PostActions):
        self.displayer.info("Found {} post actions modules".format(len(post_actions)))


