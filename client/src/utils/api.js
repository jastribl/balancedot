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

function post(path, data) {
    let options = {
        method: 'POST',
        body: data
    }
    return errorHandler(
        fetch(path, options)
    )
}

export function postJSON(path, jsonData) {
    return post(path, JSON.stringify(jsonData))
}

export function postForm(path, formData = null) {
    return post(path, formData)
}

export function postJSONWithHandling(path, setErrorMessage, setLoading) {
    setLoading(true)
    return postJSON(path)
        // .then(response => setResponse(response))
        .catch(e => setErrorMessage(e.message))
        .finally(() => setLoading(false))
}

export function get(path) {
    return errorHandler(
        fetch(path)
    )
}

export function getWithHandling(path, setResponse, setErrorMessage, setLoading) {
    setLoading(true)
    return get(path)
        .then(response => setResponse(response))
        .catch(e => setErrorMessage(e.message))
        .finally(() => setLoading(false))
}
