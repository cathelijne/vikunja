
export const saveListView = (listId, routeName) => {
	const savedListView = localStorage.getItem('listView')
	let savedListViewJson = false
	if (savedListView !== null) {
		savedListViewJson = JSON.parse(savedListView)
	}

	let listView = {}
	if(savedListViewJson) {
		listView = savedListViewJson
	}

	listView[listId] = routeName
	localStorage.setItem('listView', JSON.stringify(listView))
}

export const getListView = listId => {
	// Remove old stored settings
	const savedListView = localStorage.getItem('listView')
	if(savedListView !== null && savedListView.startsWith('list.')) {
		localStorage.removeItem('listView')
	}

	if (!savedListView) {
		return 'list.list'
	}

	const savedListViewJson = JSON.parse(savedListView)

	if(!savedListViewJson[listId]) {
		return 'list.list'
	}

	return savedListViewJson[listId]
}