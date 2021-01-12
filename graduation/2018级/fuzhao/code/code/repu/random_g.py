import random
import math
#生成随机数
def random_nr(n, num): #不重复随机数，范围，个数
	a = [0] * n
	for i in range(0, n):
		a[i] = i
	ans = random.sample(a,num)
	#print(ans)
	return ans
	
#print(math.exp(10), math.log(1))
#random_nr(10, 5)