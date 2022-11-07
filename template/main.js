

// 折线图
function lineStack(v, divId, chartName) {
    var div = document.getElementById(divId);
    var myChart = echarts.init(div, null, {
        renderer: 'canvas',
        useDirtyRect: false
    });
    var app = {};

    var option;

    option = {
        title: {
            text: chartName
        },
        tooltip: {
            trigger: 'axis',
            confine: true
        },
        legend: {
            data: v.legend,
            show: true,
            top: "3%"
        },
        grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            containLabel: true
        },
        toolbox: {
            feature: {
                saveAsImage: {}
            }
        },
        xAxis: {
            type: 'category',
            boundaryGap: false,
            data: v.xAxis
        },
        yAxis: {
            type: 'value'
        },
        series: v.series
    };


    if (option && typeof option === 'object') {
        myChart.setOption(option);
    }

    window.addEventListener('resize', myChart.resize);
}

// 饼图
function pie(v, divId, chartName) {
    var dom = document.getElementById(divId);
    var myChart = echarts.init(dom, null, {
        renderer: 'canvas',
        useDirtyRect: false
    });
    var app = {};

    var option;

    option = {
        title: {
            text: chartName,
            subtext: '',
            left: 'center'
        },
        tooltip: {
            trigger: 'item',
            formatter: "{a} {b}"
        },
        legend: {
            orient: 'vertical',
            left: 'left'
        },
        series: [
            {
                name: '',
                type: 'pie',
                radius: '61.8%',
                data: v.series,
                emphasis: {
                    itemStyle: {
                        shadowBlur: 10,
                        shadowOffsetX: 0,
                        shadowColor: 'rgba(0, 0, 0, 0.5)'
                    }
                }
            }
        ]
    };

    if (option && typeof option === 'object') {
        myChart.setOption(option);
    }

    window.addEventListener('resize', myChart.resize);
}
