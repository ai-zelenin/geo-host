ymaps.ready(['projection.wgs84Mercator', 'projection.sphericalMercator']).then(function () {
    const remoteObjectManager = new ymaps.RemoteObjectManager('/api/v1/yandex?tiles=%t&bbox=%b&zoom=%z&debug=true&clusterLevel=2', {});
    remoteObjectManager.setFilter(function (object) {
        for (let key in object.properties.options) {
            object.options[key] = object.properties.options[key]
        }
        return true;
    });
    const mapSettings = {
        center: [55.756363, 37.623270],
        zoom: 10,
        controls: ['zoomControl', 'searchControl', 'typeSelector',  'fullscreenControl', 'routeButtonControl']
    }
    const projSettings = {
        projection: ymaps.projection.sphericalMercator
    }
    const map = new ymaps.Map('map', mapSettings);
    map.geoObjects.add(remoteObjectManager);
    map.controls.get('zoomControl').options.set({ size: 'small' });

});
