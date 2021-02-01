import React from 'react'

const Modal = ({ headerText, handleClose, visible, children }) => {
    const showHideClassName = visible ? 'modal display-block' : 'modal display-none'

    const onClickHandler = (event, data) => {
        if (event.target.id === 'modal-background') {
            handleClose()
        }
    }

    return (
        <div id='modal-background' className={showHideClassName} onClick={onClickHandler}>
            <section className='modal-main'>
                <div className={showHideClassName + ' modal-content '}>
                    <div className='modal-header'>
                        <span onClick={handleClose} className='close'>&times;</span>
                        <h2>{headerText}</h2>
                    </div>
                    <div className='modal-body'>
                        {children}
                    </div>
                </div>
            </section>
        </div>
    )
}

export default Modal
