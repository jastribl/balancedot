import React from 'react'

import Table from '../common/Table'

const ExtendableTable = ({
    data,
    columns,
    customRenders,
    extraColumns,
    extraCustomRenders,
    ...props
}) => {
    columns = columns.concat(extraColumns ?? [])
    Object.assign(customRenders, extraCustomRenders ?? {})

    return <Table
        rowKey='uuid'
        rows={data ?? []}
        columns={columns}
        customRenders={customRenders}
        {...props}
    />
}

export default ExtendableTable
