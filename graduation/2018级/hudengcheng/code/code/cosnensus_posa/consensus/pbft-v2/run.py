import os
from threading import Thread

def execute_cmd_command_build_exe():
    os.system("go build -o pbft31.exe")

def execute_cmd_command_run_node(number):
    command_line = "pbft31.exe N"+str(number)
    os.system(command_line)

# 程序执行脚本启动
execute_cmd_command_build_exe()
for i in range(31):
    t = Thread(target=execute_cmd_command_run_node(i))
    t.start()


