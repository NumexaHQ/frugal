import { init } from '@rematch/core'
import createPersistPlugin, { getPersistor } from '@rematch/persist'
import React from 'react'
import ReactDOM from 'react-dom'
import { Provider } from 'react-redux'
import { PersistGate } from 'redux-persist/es/integration/react'
import storage from 'redux-persist/lib/storage'
import App from './app'
import * as models from './store/model'

const persistPlugin = createPersistPlugin({
	key: 'root',
	storage,
	version: 1,
	whitelist: ['CommonState'],
})


const store = init({
	models,
	plugins: [persistPlugin],
})

ReactDOM.render(
	<React.StrictMode>
		<Provider store={store}>
			<PersistGate persistor={getPersistor()}>
			<App />
			</PersistGate>
		</Provider>
	</React.StrictMode>,
	document.getElementById('root')
)