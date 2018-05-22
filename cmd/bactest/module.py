from typing import List
from displayer import Framework
import inspect

class Module:
    def __init__(self, name: str, directory: str, files: List[str]):
        self.name = name
        self.directory = directory
        self.files = files
        self.executables = list() # imported executables modules
        self.displayer = Framework(module_name=__class__.__name__)


    def __str__(self):
        return "Name: " + str(self.name) +  " Directory: " + str(self.directory) + " Files: " + str(self.files)

    def run(self):
        self.displayer.info("Starting tests from module: {}".format(self.name))
        for m in self.executables:
            self.__setup_module_file(m)
            self.__run_module_file_tests(m)
            self.__teardown_module_file(m)

    def __setup_module_file(self, imp_module):
        self.displayer.debug("Starting setup of: {}".format(imp_module.__name__))
        try:
            imp_module.setup()
        except Exception as e:
            # FIXME Here should be defined different format of tests errors and exceptions
            self.displayer.error("An error occured when running setup: {}".format(e))

    def __run_module_file_tests(self, imp_module):
        self.displayer.debug("Starting tests of module: {}".format(imp_module.__name__))
        test_cls = inspect.getmembers(imp_module, inspect.isclass)
        for k in test_cls:
            cls_name = k[0]
            cls_instance = k[1]()
            self.displayer.info("Running tests class: {}".format(cls_name))
            try:
                cls_instance.run()
            except Exception as e:
                # FIXME Here should be defined different format of tests errors and exceptions
                self.displayer.error("An error occured when running test: {}".format(e))
        
    def __teardown_module_file(self, imp_module):
        self.displayer.debug("Starting teardown of: {}".format(imp_module.__name__))
        try:
            imp_module.teardown()
        except Exception as e:
            # FIXME Here should be defined different format of tests errors and exceptions
            self.displayer.error("An error occured when running teardown: {}".format(e))
