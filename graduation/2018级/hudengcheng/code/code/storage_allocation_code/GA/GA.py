from deap import algorithms, creator, base, tools
import random
import math
from scipy.stats import bernoulli
from objective_function.data_availability import data_storage_availability
from objective_function.data_storage_cost import data_storage_cost
from objective_function.data_storage_time import data_storage_time
from parameters.parameters import csp_number
from normalization.normalization import normalization_with_number
import numpy as np

random.seed(42)

# 描述问题
creator.create("FitnessMin", base.Fitness, weights=(-1.0,))
creator.create("Individual", list, fitness=creator.FitnessMin)
GENE_LENGTH = csp_number

toolbox = base.Toolbox()
toolbox.register('binary', bernoulli.rvs, 0.5)  # 注册一个Binary的alias，指向scipy.stats中的bernoulli.rvs，概率为0.5
toolbox.register('individual', tools.initRepeat, creator.Individual, toolbox.binary, n=GENE_LENGTH)


# 评价函数
def eval_func(ind):
    print("Availability: ", data_storage_availability(ind), "Storage: ", data_storage_time(ind), "Cost: ", data_storage_cost(ind))
    return (normalization_with_number(-math.log(data_storage_availability(ind))) + normalization_with_number(data_storage_time(ind)) + normalization_with_number(data_storage_cost(ind))),


# 在工具箱中注册遗传算法需要的工具
toolbox.register('evaluate', eval_func)
toolbox.register('select', tools.selTournament, tournsize=3)  # 注册Tournsize为2的锦标赛选择
toolbox.register('mate', tools.cxUniform, indpb=0.8)  # 注意这里的indpb需要显示给出
toolbox.register('mutate', tools.mutFlipBit, indpb=0.05)

# 生成初始族群
N_POP = 100  # 族群中的个体数量
toolbox.register('population', tools.initRepeat, list, toolbox.individual)
pop = toolbox.population(n=N_POP)
hof = tools.HallOfFame(1)


# 注册计算过程中需要记录的数据
stats = tools.Statistics(lambda ind: ind.fitness.values)
stats.register("avg", np.mean)
stats.register("std", np.std)
stats.register("min", np.min)
stats.register("max", np.max)

# 调用DEAP内置的算法
resultPop, logbook = algorithms.eaSimple(pop, toolbox, cxpb=0.9, mutpb=0.05, ngen=200, stats=stats, halloffame=hof, verbose=True)

# 输出计算过程
logbook.header = 'gen', 'nevals', "avg", "std", 'min', "max"
print(logbook)