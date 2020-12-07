import React, { useState } from 'react'

import Spinner from './Spinner'

const Form = ({ onSubmit, disableForm, fieldInfos }) => {
    const getValidatorForFieldName = (fieldName) =>
        fieldInfos[fieldName].validate ?? (() => { return null })

    let initialState = {}
    Object.entries(fieldInfos).map(([fieldName, fieldInfo]) => {
        initialState[fieldName] = fieldInfo.initialValue ?? ''
    })
    const [formState, setFormState] = useState(initialState)
    const [validationErrors, setValidationErrors] = useState({})


    const handleFormFieldChange = (event) => {
        const fieldName = event.target.name
        const fieldValue = event.target.value
        const fieldInfo = fieldInfos[fieldName];

        if (fieldName in validationErrors) {
            const validationResult = getValidatorForFieldName(fieldName)(fieldInfo.fieldLabel, fieldValue)
            if (validationResult !== null) {
                setValidationErrors({
                    ...validationErrors,
                    [fieldName]: validationResult
                })
            } else {
                setValidationErrors({
                    ...validationErrors,
                    [fieldName]: null
                })
            }
        }
        setFormState({
            ...formState,
            [fieldName]: fieldValue
        })
    }

    const onSubmitInternal = (event) => {
        event.preventDefault()
        if (!Object.entries(fieldInfos).some(([fieldName, fieldInfo]) => {
            const validationResult = getValidatorForFieldName(fieldName)(fieldInfo.fieldLabel, formState[fieldName])
            if (validationResult === null) {
                return false
            }
            setValidationErrors({
                ...validationErrors,
                [fieldName]: validationResult
            })
            return true
        })) {
            onSubmit(formState).then(() => setFormState(initialState))
        }
    }

    return (
        <form onSubmit={onSubmitInternal} autoComplete="off" style={{ position: 'relative' }}>
            <Spinner visible={disableForm} />
            <div className="row">
                {Object.entries(fieldInfos).map(([fieldName, fieldInfo]) =>
                    <div key={fieldName} className="row">
                        <div className="col-25">
                            <label>{fieldInfo.fieldLabel}</label>
                        </div>
                        <div className="col-75">
                            <input
                                type={fieldInfo.inputType}
                                name={fieldName}
                                value={formState[fieldName]}
                                onChange={handleFormFieldChange}
                                placeholder={fieldInfo.placeholder}
                                disabled={disableForm}
                            />
                            <span style={{ float: 'right', color: 'red' }}>{validationErrors[fieldName]}</span>
                        </div>
                    </div>
                )}
            </div>
            <div className="row">
                <input type="submit" value="Add" disabled={disableForm} />
            </div>
        </form >
    )
}

export default Form
