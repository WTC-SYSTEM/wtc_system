window.onload = function () {
    //<editor-fold desc="Changeable Configuration Block">

    // the following lines will be replaced by docker/configurator, when it runs in a docker-container
    window.ui = SwaggerUIBundle({
        url: "http://www.michaeldeveloper.com/public/swagger.yaml",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIStandalonePreset
        ],
        plugins: [
            SwaggerUIBundle.plugins.DownloadUrl,
            DisableTryItOutPlugin
        ],
        layout: "StandaloneLayout"
    });

    //</editor-fold>
};

const DisableTryItOutPlugin = function () {
    return {
        statePlugins: {
            spec: {
                wrapSelectors: {
                    allowTryItOutFor: () => () => false
                }
            }
        }
    }
}
