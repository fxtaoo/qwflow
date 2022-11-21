

// 折线图
function lineStack(v, divId, chartName, downloadImg, prefix) {
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
        textStyle: {
            fontSize: 12,
            fontWeight: "normal",
        },
        legend: {
            data: v.legend,
            show: true,
            top: "5%"
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

    if (downloadImg === "true") {
        option.textStyle.fontSize = 18;
        option.textStyle.fontWeight = "bold";
    }

    if (option && typeof option === 'object') {
        myChart.setOption(option);
    }

    window.addEventListener('resize', myChart.resize);
    setTimeout(() => {
        dataUrl = myChart.getDataURL();
        imgUrl = myChart.getDataURL({
            pixelRatio: 2,
            backgroundColor: '#fff'
        });
        if (downloadImg == "true") {
            saveImg(imgUrl, prefix + "-" + divId)
        }
    }, 1000);
}

// 饼图
function pie(v, divId, chartName, downloadImg, prefix) {
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
        textStyle: {
            fontSize: 12,
            fontWeight: "normal",
        },
        legend: {
            orient: 'vertical',
            left: 'left',
            show: true,
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

    if (downloadImg === "true") {
        option.textStyle.fontSize = 18;
        option.textStyle.fontWeight = "bold";
        option.legend.show = false;
    }

    if (option && typeof option === 'object') {
        myChart.setOption(option);
    }

    window.addEventListener('resize', myChart.resize);
    setTimeout(() => {
        dataUrl = myChart.getDataURL();
        imgUrl = myChart.getDataURL({
            pixelRatio: 2,
            backgroundColor: '#fff'
        });
        if (downloadImg == "true") {
            saveImg(imgUrl, prefix + "-" + divId)
        }
    }, 1000);
}

function saveImg(imgUrl, divId) {
    const a = document.createElement('a')
    a.href = imgUrl
    a.setAttribute('download', divId)
    a.click()
}
