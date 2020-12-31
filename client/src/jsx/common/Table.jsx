import React, { useState } from 'react'

const Table = ({ rowKey, columns, rows, customRenders, initialSortColumn, customSortComparators }) => {
    customRenders ??= {}
    customSortComparators ??= {}

    if (!rows) {
        return <div />
    }

    const [sortColumn, setSortColumn] = useState(initialSortColumn)
    const [sortInverse, setSortInverse] = useState(false)

    const onHeaderClick = (header_name) => {
        if (sortInverse) {
            setSortColumn(null)
            setSortInverse(false)
        } else if (sortColumn === header_name) {
            setSortInverse(!sortInverse)
        } else {
            setSortColumn(header_name)
            setSortInverse(false)
        }
    }

    let toRender = rows.slice()
    if (sortColumn) {
        toRender.sort(customSortComparators[sortColumn] ?? ((a, b) => {
            if (a[sortColumn] < b[sortColumn]) {
                return -1
            } else if (a[sortColumn] > b[sortColumn]) {
                return 1
            } else {
                return 0
            }
        }))
    }
    if (sortInverse) {
        toRender.reverse()
    }
    return (
        <div>
            <table className='styled-table'>
                <thead>
                    <tr>
                        {Object.keys(columns).map((key) =>
                            <th
                                key={key}
                                onClick={() => onHeaderClick(key)}
                            >{columns[key]}{(key === sortColumn ? (sortInverse ? " ↑" : " ↓") : "")}</th>
                        )}
                    </tr>
                </thead>
                <tbody>
                    {toRender.map((row, i) =>
                        <tr key={row[rowKey]}>{
                            Object.keys(columns).map((key) =>
                                <td key={key}>{
                                    key in customRenders ? customRenders[key](row) : row[key]
                                }</td>
                            )
                        }</tr>
                    )}
                </tbody>
            </table>
        </div>
    )
}

export default Table