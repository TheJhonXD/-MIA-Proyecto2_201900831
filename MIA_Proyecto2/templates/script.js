$(document).ready(function(){
    inputPath = document.querySelector('#inputPath');
    myfile = document.querySelector('#myfile');
    editor = document.querySelector('#editor');
    result = document.querySelector('#result');
    btnMenuMain = document.querySelector('#btn-menu-main')
    btnMenuLogin = document.querySelector('#btn-menu-login')
    btnMenuReport = document.querySelector('#btn-menu-reports')
    btnLogin = document.querySelector('#btn-ingresar')
    btnEjecutar = document.querySelector('#execute')
    btnGraph = document.querySelector('#btn-graficar')
    selectGraph = document.querySelector('#graphs')


    function readFile(file, callback){
        const reader = new FileReader();
        reader.addEventListener('load', (event) => {
            const content = event.target.result;
            callback(content)
        });
          
        reader.readAsText(file);
    }

    myfile.addEventListener('change', (e) => {
        inputPath.value = e.target.files[0].name;
        readFile(e.target.files[0], (content) => {
            editor.value = content;
        })
    });

    btnMenuMain.addEventListener('click', (e) => {
        $('#main').siblings().slideUp(500);
        $('#main').slideDown(500);
    });

    btnMenuLogin.addEventListener('click', (e) => {
        $('#login-page').siblings().slideUp(500);
        $('#login-page').slideDown(500);
    });

    btnMenuReport.addEventListener('click', (e) => {
        $('#report-page').siblings().slideUp(500);
        $('#report-page').slideDown(500);
    });

    btnEjecutar.addEventListener('click', async () => {
        let data = {
            "message": `${editor.value}`
        }

        await fetch('http://localhost:3000/text', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Error al hacer la solicitud');
            }
            return response.json();
        })
        .then(data => {
            result.value = data.message;
            console.log("OK");
        })
        .catch(error => console.error(error));
    });

    btnLogin.addEventListener('click', (e) => {
        e.preventDefault();
        sessionStorage.setItem('user', $('#user-name').val());
        sessionStorage.setItem('id', $('#id-part').val());
        alert(`"${sessionStorage.getItem('user')} a iniciado sesion."`)
    });


    btnGraph.addEventListener('click', async () => {
        let data = {
            "name": selectGraph.value,
            "id": sessionStorage.getItem('id')
        }
        await fetch('http://localhost:3000/graph', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Error al hacer la solicitud');
            }
            return response.json();
        })
        .then(data => {
            console.log(data.message);
            d3.select("#thumb").graphviz({useWorker: false})
            .fade(false)
            .renderDot(data.message);
            console.log("OK");
        })
        .catch(error => console.error(error));

    });
});