num=xlsread('size_height.xlsx')

x = num(:,1)
y = num(:,2)
b = num(:,3)
d = num(:,4)

plot(x,y,'-*',x,b,'-s',x,d,'-o'),axis ( [0 2100 0 50] ),xlabel('���ݵ�����'), ylabel('�洢�ռ�����(kb)'),legend('DFHMT-11','DFHMT-4','DFHMT-8')