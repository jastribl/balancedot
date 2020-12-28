function errorHandler(pm) {
    return pm
        .then(response => {
            if (response.ok) {
                return response.json()
            }
            return response.json().then(e => {
                e.status = response.status
                throw e
            })
        })
}

function post(path, data, redirect = null) {
    let options = {
        method: 'POST',
        body: data
    }
    if (redirect !== null) {
        options['redirect'] = redirect
    }
    return errorHandler(
        fetch(path, options)
    )
}

export function postJSON(path, jsonData, redirect = null) {
    return post(path, JSON.stringify(jsonData), redirect)
}

export function postForm(path, formData = null, redirect = null) {
    return post(path, formData, redirect)
}

export function get(path) {
    return errorHandler(
        fetch(path)
    )
}
