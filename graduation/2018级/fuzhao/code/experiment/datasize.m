num=xlsread('datasize.xlsx')
x = num(:,1)
y = num(:,2)
z = num(:,3)
a = num(:,4)
c = num(:,5)



plot(x,y,'-*',x,z,'-s',x,a,'-^',x,c,'-d'),axis ( [0 2100 0 250] ),xlabel('���ݵ�����'), ylabel('�洢�ռ�����(kb)'),legend('DFHMT','B-DIS','B-DAM','BB-DIS')