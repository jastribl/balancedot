import React, { useState } from 'react'

import { defaultSort } from '../../utils/sorting'
import { snakeToSentenceCase } from '../../utils/strings'

const Table = ({
    rowKey,
    columns,
    rows,
    customRenders,
    initialSortColumn,
    initialSortInverse,
    customSortComparators,
    hideFilters,
    customSortFunctionOverride,
}) => {
    customRenders ??= {}
    customSortComparators ??= {}
    hideFilters ??= false

    if (!rows) {
        return <div />
    }

    const [sortColumn, setSortColumn] = useState(initialSortColumn)
    const [sortInverse, setSortInverse] = useState(initialSortInverse ?? false)

    let initialFilters = {}
    columns.map(key => {
        initialFilters[key] = ''
    })
    const [filters, setFilters] = useState(initialFilters)

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

    const handleFilterChange = (event) => {
        const filterKey = event.target.name
        const filterValue = event.target.value
        setFilters({
            ...filters,
            [filterKey]: filterValue,
        })
    }

    let toRender = rows.slice()

    toRender = toRender.filter(row => {
        return !Object.keys(filters).some(filterKey => {
            const rawValue = filterKey in customRenders ? customRenders[filterKey](row) : row[filterKey]
            const displayValue = ('' + rawValue).toLowerCase()
            let searchTerm = filters[filterKey].toLowerCase()
            if (searchTerm.length === 0) {
                return
            }
            let isInverse = searchTerm.length > 1 && searchTerm.startsWith('!')
            if (isInverse) {
                searchTerm = searchTerm.substring(1)
            }
            const contains = displayValue.includes(searchTerm)
            return isInverse ? contains : !contains
        })
    })

    if (sortColumn) {
        toRender.sort((a, b) => (customSortComparators[sortColumn] ?? defaultSort)(
            a[sortColumn],
            b[sortColumn],
        ))
        if (sortInverse) {
            toRender.reverse()
        }
    } else if (customSortFunctionOverride) {
        toRender.sort(customSortFunctionOverride)
    }

    let filterDiv = null
    if (!hideFilters) {
        filterDiv = <tr>
            {columns.map(key =>
                <td key={key} >
                    <input
                        type={'text'}
                        name={key}
                        value={filters[key]}
                        onChange={handleFilterChange}
                        placeholder={'Filter for ' + snakeToSentenceCase(key)}
                    />
                </td>
            )}
        </tr>
    }

    return (
        <div>
            <table className='styled-table'>
                <thead>
                    {filterDiv}
                    <tr>
                        {columns.map(key =>
                            <th
                                key={key}
                                onClick={() => onHeaderClick(key)}
                            >{snakeToSentenceCase(key)}{(key === sortColumn ? (sortInverse ? ' ↑' : ' ↓') : '')}</th>
                        )}
                    </tr>
                </thead>
                <tbody>
                    {toRender.map((row, _i) =>
                        <tr key={row[rowKey]}>{
                            columns.map(key =>
                                <td key={key}>{
                                    key in customRenders ? customRenders[key](row) : row[key]
                                }</td>
                            )
                        }</tr>
                    )}
                </tbody>
            </table>
        </div >
    )
}

export default Table