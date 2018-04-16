# BACTEST

It's a common tool/framework used for testing backer solution. The design is focused to provide a possibility
to use it also in other projects.


## Assumptions in nutshell
1. As tests source uses pointed directory, under this directory each dir, which has "test_*" prefix is taken
into account as a testcase 
2. Under "test_*" directory the structure of python source files must be like that:
⋅⋅* t_setup.py
⋅⋅* t_test1.py
⋅⋅* t_test2.py
⋅⋅* t_teardown.py
From this files every method from class will be executed.

3. Methods from dir "pre_*" are always executed before all tests in random order
4. Methods from dir "post_*" are always executed after all tests in random order  
5. Thinking :)