num=xlsread('datasize.xlsx')
x = num(:,1)
y = num(:,2)
z = num(:,3)
a = num(:,4)
c = num(:,5)



plot(x,y,'-*',x,z,'-s',x,a,'-^',x,c,'-d'),axis ( [0 2100 0 250] ),xlabel('数据的数量'), ylabel('存储空间消耗(kb)'),legend('DFHMT','B-DIS','B-DAM','BB-DIS')