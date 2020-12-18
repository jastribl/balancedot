
function post (path, data) {
    return fetch(path, {
        method: 'POST',
        body: data
    })
        .then(response => {
            if (response.ok) {
                return response.json()
            }
            return response.json().then(e => {
                throw {
                    status: response.status,
                    message: e.message
                }
            })
        })
}

export function postJSON(path, jsonData) {
    return post(path, JSON.stringify(jsonData))
}

export function postForm(path, formData) {
    return post(path, formData)
}

export function get(path) {
    return fetch(path)
        .then((res) => res.json())
}
