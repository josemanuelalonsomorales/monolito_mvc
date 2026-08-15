async function cargarUsuarios() {
    try {
        const respuestaServidor = await fetch('/api/usuarios');
        if (!respuestaServidor.ok) throw new Error('Error al consultar los usuarios');

        const listaUsuarios = await respuestaServidor.json();
        let contenidoHTML = '<ul>';
        
        listaUsuarios.forEach(usuarioActual => {
            contenidoHTML += `<li><b>[ID: ${usuarioActual.id}]</b> ${usuarioActual.nombre} - Saldo: $${usuarioActual.saldo.toFixed(2)}</li>`;
        });
        contenidoHTML += '</ul>';

        document.getElementById('contenedorListaUsuarios').innerHTML = contenidoHTML;
    } catch (errorCapturado) {
        document.getElementById('contenedorListaUsuarios').innerHTML = '<p style="color:red">Error al cargar los usuarios.</p>';
        console.error(errorCapturado);
    }
}

document.getElementById('formularioTransferencia').addEventListener('submit', async (eventoFormulario) => {
    eventoFormulario.preventDefault();

    const datosTransferencia = {
        id_origen: parseInt(document.getElementById('identificadorUsuarioOrigen').value),
        id_destino: parseInt(document.getElementById('identificadorUsuarioDestino').value),
        monto: parseFloat(document.getElementById('montoATransferir').value)
    };

    try {
        const respuestaServidor = await fetch('/api/transferir', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(datosTransferencia)
        });

        if (respuestaServidor.ok) {
            alert('Transferencia realizada con éxito.');
            cargarUsuarios();
        } else {
            const textoMensajeError = await respuestaServidor.text();
            alert('Error: ' + textoMensajeError);
        }
    } catch (errorCapturado) {
        alert('Error conectando con el servidor en Go.');
        console.error(errorCapturado);
    }
});

cargarUsuarios();