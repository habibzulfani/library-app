package queries

const CombinedResourcesQuery = `
    SELECT 'book' as type, id, title, author FROM books
    UNION ALL
    SELECT 'paper' as type, id, title, author FROM papers
`
