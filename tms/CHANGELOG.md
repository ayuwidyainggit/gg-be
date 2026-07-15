# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## 0.3.0 (2024-10-09)


### Features

* add shipment report detail and reject ([f9b7dec](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/f9b7deca7885acf177372974ccbd1b60e33810a2))
* add shipment report summary ([301ebd3](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/301ebd3a36b251bd6cc1bdcbc35f5dda5baf332a))
* add start time and end time ([8739d08](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/8739d08b1eb97a1d5517ed73d44fe9beffa76942))
* add vat and vat_value list invoices ([0116323](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/0116323cb77633e5a9dc0b43015530601ee939f6))
* adding docker file ([b2fc70b](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/b2fc70b4ea9c0374ad55ef4f198f2393e93a8cbd))
* adding github action ([976c82e](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/976c82ea79a8e38b360ec9ca419074cccd81ec0e))
* cancel reject & filter driver_id in shipments ([5b34979](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/5b3497938b312d7cca2f5fae7f772196cd726c8b))
* create shipment auto, fix total product in outlet, get the vehicle and it doesn't show up the existing vehicle_id, get data shipment_auto via sendpick(third party), add the shipment_no column looks for data that has been reject all ([655214c](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/655214cd3580595fea2dc11a7dd73502c562dc3e))
* delete bulk shipment and get vehicle filter delivery_date ([908efd1](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/908efd1d94ad09227420a0ac78da62b9496185f3))
* driver report mobile ([3f3748e](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/3f3748e8f7b0287d4c2e16c7135e181b58e11945))
* new request skip and response todo mobile ([6e73221](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/6e73221feb308102ce20787a49c1355b6f3e35aa))
* **pickup:** implement api pickup ([695941b](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/695941ba2b016d5ce5f59e11bd2c629927a2d0ea))
* start visit & arrive mobile ([f9bec92](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/f9bec92adf578078f9fe9bba43c32ea2a2cd138b))
* unload, arrive, leave, reason, reject all & partial ([a84c97d](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/a84c97db9c458e0aa7cf6bad2795d29c49fe35a0))


### Bug Fixes

* Add column shipping_type & Add param for reject all & partial ([81c4a78](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/81c4a78ee4986deab124d6c0b0860c1d35ed5cba))
* add reason_id, reason_name and current_time in reject all mobile ([4282e7b](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/4282e7b73c79d3f47a33faaee87f41cba5de8481))
* Add validation status unauthorized or empty token ([893240b](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/893240b17063bea8ef6d25e5775052c71cf7d4ae))
* adding credential for staging ([288d316](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/288d316fcdeb33447f96e9efb9d0539e80f08560))
* arrive,skip,unload,start,end & filter shipments by delivery_date and cust_id ([719d089](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/719d0895c13d61f7312c178e5734aea57b431d5b))
* count trip in visit summary mobile ([785251c](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/785251c907e3575fd8a5b3329846fb701e4849ba))
* create manual shipment(web), get outlets params shipment_no(mobile), delete shipment(web), get shipment(web) ([91d377f](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/91d377ffa294eb45e01616e8be9c37de799adf89))
* create shipment auto when invoice duplicate and list reject partial and send reject partial ([aa8b421](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/aa8b4212df6e977c06ee4d6a5b8adeaf8ba5eb12))
* create shipment auto when vehicle less than total shipment_no ([804aad1](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/804aad149ec5f248583f14a7e99eece0d4b15625))
* docker file ([479b2fd](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/479b2fd438647036c68d05012743e11740fec3b6))
* endpoint hit for filter invoice ([4459cef](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/4459cef1c1bbea0de23b09578b38112c8de870f1))
* endpoint hit for filter invoice edpoint kong port changed to 9004 ([a7038d6](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/a7038d6bcd9fb83b32629b27b1ef025c2e6a0262))
* get all invoice ([8d3afdc](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/8d3afdcc0205467a84503f36c8a13dd20e6e40df))
* get invoices ([1b2b249](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/1b2b249a3dc02a21f66ec7b18f7f42acad14f382))
* get invoices ([46224f2](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/46224f276d52148f9568710d9f88ec2fdef2c81f))
* get invoices ([25d423a](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/25d423a9ae300822ac9c3f417a02dcda90e0c810))
* get product add key id, reject partial body aoa(array of object), reject all rename the key to id ([ecfdbba](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/ecfdbba9652ab70265f8e8ea2616fb28bf3e0d87))
* handle query update and delete when rows affected equals 0 ([df13309](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/df1330981ad065bb75e8974cf325f5c0934f664d))
* handling error thirdparty in create shipment auto ([0b49b95](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/0b49b95dcffe058e72570dcd632c0d04af7937d0))
* implement params status for get product mobile and add attr total_product_delivery and total_product_pickup ([5f5ce0a](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/5f5ce0a2a2bde308ba1c33c4933fdd44c9e1a75a))
* logic query driver report mobile ([81578d7](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/81578d7e486902d6ff776b9af612f4928f8ae628))
* **pickup:** remove auth in pickup and add response attr in product ([8f8fb5f](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/8f8fb5f2376d1f63b25e0cf22715a884fa4da5ea))
* **pickup:** update structure request pickup partial and skip pickup ([4b731c9](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/4b731c9e6b026664ad572673d7e1dff85fa4e2ad))
* pipeline to add --build ([91bcd86](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/91bcd86e1704d054334e7548d854f939f3bc8123))
* reject cancel, create shipment auto duplicate vehicle when vehicle_id duplicate and missing key response json in list reject partial ([87f9d5b](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/87f9d5b2f6523b1630b255aff50928b1393e262d))
* Reject Partial and product_status ([caaa4ee](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/caaa4ee979b794f64c0a7069bdda9ca4d5e7473e))
* remove /swagger ([5b92c7f](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/5b92c7f22d8f99f8f220785c849505d7ab53c8ff))
* rename invoice_no type int to string ([cacfbd9](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/cacfbd932e237bc03a59a73f309475ef472e9336))
* rewrite codebase service and router, add is_active is daily activity mobile ([f6f4d6b](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/f6f4d6b82caa98d7bcf309049025a22eda085271))
* rewrite url mobile & web, create shipment web, summary, cancel ([7d566c2](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/7d566c2da7145411df4d933a52cd32850662207e))
* submit shipment preview in web ([1cd2aa1](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/1cd2aa1c0eb5996723e025fdfe5c76bb72101489))
* submit start & finish, rewrite codebase, change type column ([fd85e15](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/fd85e15906b8dc6153bb92edf3870addc74cc163))
* todo list in mobile ([8b53681](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/8b53681790f2c651853e85e41bd41375408ed986))
* update status shipment, shipment daily acitivity for the current date, logic summary driver, add payload shipment_no in start, end, leave,arrive, unload and todo ([c59d920](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/c59d920bd817c232ad8d4e46fc69caec3439c98d))
* visit leave and skip so that their status values become Finished and Skipped ([00664e1](https://github.com/repodigithub/ScyllaX-TMS-BE/commit/00664e190a6fa2d28ca2d32fa869e6b464ab84e3))
