const editFetch = (form, inputName, inputDate, inputDesc, tags) => {
    form.addEventListener('submit', (e) => {
        e.preventDefault();
        const testName = inputName.value.length === 0;
        const testDate = inputDate.value.toString().length === 0;
        const testDesc = inputDesc.value.length === 0;

        if (!testName) {
            inputName.className = 'form-field-edit-novalid text-without-icon';
        }

        if (!testDate) {
            inputDate.className = 'form-field-edit-novalid text-with-icon';
        }

        if (!testDesc) {
            inputDesc.className = 'form-field-edit-novalid';
        }

        if (!testName || !testDate || !testDesc) {
            return;
        }

        const name = inputName.value.trim();
        const date = inputDate.value.trim();
        const desc = inputDesc.value.trim();
        const requestOptions = {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                'name': name,
                'age': date,
                'description': desc,
                'tags': tags
            })
        };
        fetch(`${localAddress}/api/v1/edit`, requestOptions)
            .then(response =>
                response.json().then(data => ({
                    data: data,
                    status: response.status
                })).then(res => {
                    if (res.status === 200 && res.data.status === 200) {
                        clearRoot();
                        setUserProfile(res.data.body);
                        renderEdit();
                    } else if (res.data.status === 404) {
                        const userNotFound = document.createElement('span')
                        userNotFound.textContent = 'Неправильный ввод'
                        userNotFound.style.marginTop = "10px"
                        form.appendChild(userNotFound)
                    }
                })).catch((error) => console.log(error));
    })
}