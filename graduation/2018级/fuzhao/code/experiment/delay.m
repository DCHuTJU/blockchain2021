





num=xlsread('delay.xlsx')
a = num(1,:)
poro = num(2,:)
pow = num(3,:)


figure;
hold on;
h1 = bar(a, [pow',poro']);
set(gca,'xtick',[1 2 3 4 5 6 7 8 9 10]);
set(gca,'xticklabel',{'0-0.5','0.5-1','1-1.5','1.5-2','2-2.5','2.5-3','3-3.5','3.5-4','4-4.5','4.5-5'}),xlabel('信誉值偏移量'), ylabel('延迟(块高度差)');
legend([h1(1), h1(2)],'POW','PORO');