import React, { useState } from 'react'

import Spinner from './Spinner'
import ErrorRow from './ErrorRow'

const Form = ({ onSubmit, fieldInfos }) => {
    const getValidatorForFieldName = (fieldName) =>
        fieldInfos[fieldName].validate ?? (() => { return null })

    let initialValues = {}
    Object.entries(fieldInfos).map(([fieldName, fieldInfo]) => {
        initialValues[fieldName] = fieldInfo.initialValue ?? ''
    })
    const [formState, setFormState] = useState(initialValues)
    const [formValues, setFormValues] = useState(initialValues)
    const [validationErrors, setValidationErrors] = useState({})
    const [errorMessage, setErrorMessage] = useState(null)
    const [isSubmitting, setIsSubmitting] = useState(false)


    const handleFormFieldChange = (event) => {
        const fieldName = event.target.name
        const fieldInfo = fieldInfos[fieldName]
        const fieldValue = event.target.value
        setFormValues({
            ...formValues,
            [fieldName]: fieldValue
        })
        let formValue = fieldValue;
        if (fieldInfo.inputType === 'file') {
            formValue = event.target.files[0]
        }

        if (fieldName in validationErrors) {
            const validationResult = getValidatorForFieldName(fieldName)(fieldInfo.fieldLabel, formValue)
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
            [fieldName]: formValue
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
            setErrorMessage(null)
            setIsSubmitting(true)
            onSubmit(formState)
                .then(() => {
                    setFormState(initialValues)
                })
                .catch((e) => {
                    setErrorMessage(e)
                })
                .finally(() => {
                    setIsSubmitting(false)
                })
        }
    }

    return (
        <form onSubmit={onSubmitInternal} autoComplete="off" style={{ position: 'relative' }}>
            <Spinner visible={isSubmitting} />
            <ErrorRow message={errorMessage} />
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
                                value={formValues[fieldName]}
                                onChange={handleFormFieldChange}
                                placeholder={fieldInfo.placeholder}
                                disabled={isSubmitting}
                            />
                            <span style={{ float: 'right', color: 'red' }}>{validationErrors[fieldName]}</span>
                        </div>
                    </div>
                )}
            </div>
            <div className="row">
                <input type="submit" value="Add" disabled={isSubmitting} />
            </div>
        </form >
    )
}

export default Form
