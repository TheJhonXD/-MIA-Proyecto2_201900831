$(document).ready(function(){
    let inputPath = document.querySelector('#inputPath');
    let myfile = document.querySelector('#myfile');
    let editor = document.querySelector('#editor');
    let result = document.querySelector('#result');
    let btnMenuMain = document.querySelector('#btn-menu-main')
    let btnMenuLogin = document.querySelector('#btn-menu-login')
    let btnMenuReport = document.querySelector('#btn-menu-reports')
    let btnLogin = document.querySelector('#btn-ingresar')
    let btnLogout = document.querySelector('#btn-logout')
    let btnEjecutar = document.querySelector('#execute')
    let btnGraph = document.querySelector('#btn-graficar')
    let selectGraph = document.querySelector('#graphs')
    let resultFile = document.querySelector('#reporteFile')
    let hostAPI = "http://localhost:3000"



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

        await fetch(`${hostAPI}/text`, {
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

    btnLogin.addEventListener('click', async (e) => {
        e.preventDefault();

        let data = {
            "user": $('#user-name').val(),
            "password": $('#pass').val(),
            "id": $('#id-part').val()
        }

        await fetch(`${hostAPI}/login`, {
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
            // console.log(data)
            sessionStorage.setItem('user', $('#user-name').val());
            sessionStorage.setItem('id', $('#id-part').val());
            alert(`${data.message}`)
        })
        .catch(error => console.error(error));
    });

    btnLogout.addEventListener('click', async () =>{
        await fetch(`${hostAPI}/logout`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Error al hacer la solicitud');
            }
            return response.json();
        })
        .then(data => {
            sessionStorage.clear();
            alert(`${data.message}`)
        })
        .catch(error => console.error(error));
    })


    btnGraph.addEventListener('click', async () => {
        let data = {
            "name": selectGraph.value,
            "id": sessionStorage.getItem('id')
        }
        await fetch(`${hostAPI}/graph`, {
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
            if (selectGraph.value != "graph-file"){
                d3.select("#thumb").graphviz({useWorker: false})
                .fade(false)
                .renderDot(data.message);
            }else{
                $('#rFile').html(data.message)
            }
            console.log("OK");
        })
        .catch(error => console.error(error));

    });
});