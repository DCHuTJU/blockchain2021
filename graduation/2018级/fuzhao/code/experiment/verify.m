



num=xlsread('verify.xlsx')
x = num(:,1)
y = num(:,2)
z = num(:,3)
h = num(:,4)
b = num(:,5)
e = num(:,6)

plot(x,y,'-*',x,h,'-o',x,b,'-^',x,e,'-d'),axis ( [0 1100 0 500] ),xlabel('数据的数量'), ylabel('平均验证时间(ms)'),legend('DFHMT-11','B-DIS','DFHMT-4','DFHMT-8')