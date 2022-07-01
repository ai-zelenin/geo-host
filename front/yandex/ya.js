ymaps.ready(['projection.wgs84Mercator', 'projection.sphericalMercator']).then(function () {
    const url = new URL(window.location.href);
    const params = new URLSearchParams(url.search);
    let debug = params.get("debug");
    let clusterDepth = params.get("clusterDepth");
    if (!debug){
        debug = "false"
    }
    if (!clusterDepth){
        clusterDepth = "1"
    }
    console.log(window.location);
    const remoteObjectManager = new ymaps.RemoteObjectManager(`/api/v1/yandex?tiles=%t&zoom=%z&debug=${debug}&clusterDepth=${clusterDepth}`, {
        "paddingTemplate": "cb_%t_%z"
    });
    remoteObjectManager.setFilter(function (object) {
        if (object.properties.iconContent !== "") {
            console.log(object.properties.iconContent)
        }
        for (let key in object.properties.options) {
            object.options[key] = object.properties.options[key]
        }
        return true;
    });
    const mapSettings = {
        center: [55.756363, 37.623270],
        zoom: 10,
        controls: ['zoomControl', 'searchControl', 'typeSelector', 'fullscreenControl', 'routeButtonControl']
    }
    // const projSettings = {
    //     projection: ymaps.projection.sphericalMercator
    // }
    const map = new ymaps.Map('map', mapSettings);
    map.geoObjects.add(remoteObjectManager);
    map.controls.get('zoomControl').options.set({size: 'small'});

});
