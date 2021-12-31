import React, { useState } from 'react'

import { postForm } from '../../utils/api'
import Form from '../common/Form'
import LoaderComponent from '../common/LoaderComponent'
import Modal from '../common/Modal'
import CardActivitiesTable from '../tables/CardActivitiesTable'

const CardActivitiesPage = ({ match }) => {
    const cardUUID = match.params.cardUUID

    const [card, setCard] = useState(null)
    const [modalVisible, setShowModal] = useState(false)
    const [uploading, setUploading] = useState(false)

    const showModal = () => { setShowModal(true) }
    const hideModal = () => { setShowModal(false) }

    const handleActivityUpload = (activityData) => {
        setUploading(true)
        let formData = new FormData()
        for (let i = 0; i < activityData['files'].length; i++) {
            formData.append(`file${i}`, activityData['files'][i])
        }
        return postForm(`/api/cards/${cardUUID}/activities`, formData)
            .then(() => {
                hideModal()
            })
            .catch(e => { throw e.message })
            .finally(() => setUploading(false))
    }

    return (
        <div>
            <h1>Card Activities for {card ? (`${card.last_four} (${card.description})`) : null}</h1>
            <input type='button' onClick={showModal} value='Upload Activities' style={{ marginBottom: 25 + 'px' }} />
            <LoaderComponent
                path={`/api/cards/${cardUUID}`}
                parentLoading={uploading}
                setData={setCard}

            />
            <CardActivitiesTable data={card?.activities ?? []} />
            <Modal headerText='Activity Upload' visible={modalVisible} handleClose={hideModal}>
                <Form
                    onSubmit={handleActivityUpload}
                    fieldInfos={{
                        files: {
                            inputType: 'file',
                            multiple: true,
                        },
                    }}
                />
            </Modal>
        </div>
    )
}

export default CardActivitiesPage
