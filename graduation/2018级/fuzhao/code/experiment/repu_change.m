num=xlsread('repu_change.xlsx')

a = num(:,1)
c = num(:,2)
cc = num(:,3)



plot(a,c,'-*',a,cc,'-s'),set(gca,'ylim',[0,25]),xlabel('����'), ylabel('��������ֵ'),legend('BTIC','TDMS')