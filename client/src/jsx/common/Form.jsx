import React, { useState } from 'react'

import { snakeToSentenceCase } from '../../utils/strings'
import ErrorRow from './ErrorRow'
import Spinner from './Spinner'

const Form = ({ onSubmit, fieldInfos }) => {
    const getValidatorForFieldName = (fieldName) =>
        fieldInfos[fieldName].validate ?? (() => { return null })

    let initialValues = {}
    Object.entries(fieldInfos).map(([fieldName, fieldInfo]) => {
        initialValues[fieldName] = fieldInfo.initialValue ?? ''
        if (
            fieldInfo.inputType === 'select' &&
            initialValues[fieldName] === '' &&
            fieldInfo.selectOptions.length > 0) {
            initialValues[fieldName] = fieldInfo.selectOptions[0]
        }

        // Default labels
        fieldInfo.fieldLabel ??= snakeToSentenceCase(fieldName)

        // Default placeholders
        fieldInfo.placeholder ??= snakeToSentenceCase(fieldName) + "..."
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
        let formValue = fieldValue
        if (fieldInfo.inputType === 'file') {
            if (fieldInfo.multiple) {
                formValue = event.target.files
            } else {
                formValue = event.target.files[0]
            }
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
                    setFormValues(initialValues)
                })
                .catch((e) => {
                    setErrorMessage(e)
                })
                .finally(() => {
                    setIsSubmitting(false)
                })
        }
    }

    const getInputSection = (fieldName, fieldInfo) => {
        if (fieldInfo.inputType === 'select') {
            return (
                <select
                    name={fieldName}
                    value={formValues[fieldName]}
                    onChange={handleFormFieldChange}
                    disabled={isSubmitting}
                >
                    {fieldInfo.selectOptions.map(option => {
                        return (
                            <option
                                key={option}
                                value={option}
                            >{snakeToSentenceCase(option)}</option>
                        )
                    })}
                </select>)
        } else if (fieldInfo.inputType === 'textarea') { // todo: look into this (is this used)
            return <textarea
                name={fieldName}
                value={formValues[fieldName]}
                onChange={handleFormFieldChange}
                placeholder={fieldInfo.placeholder}
                disabled={isSubmitting}
                rows="10"
            />
        } else {
            return <input
                type={fieldInfo.inputType}
                name={fieldName}
                value={formValues[fieldName]}
                onChange={handleFormFieldChange}
                placeholder={fieldInfo.placeholder}
                disabled={isSubmitting}
                multiple={fieldInfo.multiple ? "multiple" : null}
            />
        }
    }

    return (
        <form onSubmit={onSubmitInternal} autoComplete='off' style={{ position: 'relative' }}>
            <Spinner visible={isSubmitting} />
            <ErrorRow message={errorMessage} />
            <div className='row'>
                {Object.entries(fieldInfos).map(([fieldName, fieldInfo]) =>
                    <div key={fieldName} className='row'>
                        <div className='col-25'>
                            <label>{fieldInfo.fieldLabel}</label>
                        </div>
                        <div className='col-75'>
                            {getInputSection(fieldName, fieldInfo)}
                            <span style={{ float: 'right', color: 'red' }}>{validationErrors[fieldName]}</span>
                        </div>
                    </div>
                )}
            </div>
            <div className='row'>
                <input type='submit' value='Add' disabled={isSubmitting} />
            </div>
        </form >
    )
}

export default Form
