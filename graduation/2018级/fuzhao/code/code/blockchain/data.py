import openpyxl
import random

data = openpyxl.load_workbook('offset.xlsx')
offset = data.worksheets[0]

num = 128 * 28

for i in range(0, num):
	offset.cell(i + 1, 2).value = random.uniform(0,5)

data.save('offset.xlsx')