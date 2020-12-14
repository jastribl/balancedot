
export function post(path, data) {
    return fetch(path, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
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

export function get(path) {
    return fetch(path)
        .then((res) => res.json())
}
