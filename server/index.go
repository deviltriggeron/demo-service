package server

import "net/http"

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Orders</title>
			<style>
				body {
					font-family: Arial, sans-serif;
				}
				.form-container {
					display: grid;
					grid-template-columns: repeat(4, 1fr);
					gap: 20px;
					align-items: start;
				}
				.form-block {
					border: 1px solid #ccc;
					padding: 15px;
					border-radius: 8px;
					background: #f9f9f9;
				}
				label {
					font-weight: bold;
				}
				input {
					width: 100%;
					margin-bottom: 10px;
					padding: 5px;
				}
				button {
					margin-top: 10px;
					padding: 10px 20px;
					border: none;
					background: #4CAF50;
					color: white;
					border-radius: 5px;
					cursor: pointer;
				}
				pre {
					margin-top: 20px;
					background: #eee;
					padding: 10px;
					border-radius: 5px;
				}
			</style>
		</head>
		<body>
			<h1>Find Order</h1>
			<input type="text" id="orderId" placeholder="Enter Order UID"/>
			<button onclick="getOrder()">Find</button>
			<button onclick="getAllOrders()">Get all orders</button>
			<pre id="result"></pre>

			<h2>Add new order</h2>
			<form id="orderForm" onsubmit="insert(); return false;">
				<div class="form-container">
					<div class="form-block">
						<h3>Orders</h3>
						<label>Order UID:</label><br/>
						<input type="text" id="orderUID" required /><br/>
						<label>Track number:</label><br/>
						<input type="text" id="trackNumber" required /><br/>
						<label>Entry:</label><br/>
						<input type="text" id="entry" required /><br/>
						<label>Locale:</label><br/>
						<input type="text" id="locale" required /><br/>
						<label>Internal signature:</label><br/>
						<input type="text" id="internal_signature" required /><br/>
						<label>Customer ID:</label><br/>
						<input type="text" id="customer_id" required /><br/>
						<label>Delivery service:</label><br/>
						<input type="text" id="delivery_service" required /><br/>
						<label>Shardkey:</label><br/>
						<input type="text" id="shardkey" required /><br/>
						<label>Sm ID:</label><br/>
						<input type="number" id="sm_id" required /><br/>
						<label>Date created:</label><br/>
						<input type="text" id="date_created" required /><br/>
						<label>Oof shard:</label><br/>
						<input type="text" id="oof_shard" required /><br/>
					</div>

					<div class="form-block">
						<h3>Delivery</h3>
						<label>Name:</label><br/>
						<input type="text" id="delivery_name" /><br/>
						<label>Phone:</label><br/>
						<input type="text" id="delivery_phone" /><br/>
						<label>Zip:</label><br/>
						<input type="text" id="delivery_zip" /><br/>
						<label>City:</label><br/>
						<input type="text" id="delivery_city" /><br/>
						<label>Address:</label><br/>
						<input type="text" id="delivery_address" /><br/>
						<label>Region:</label><br/>
						<input type="text" id="delivery_region" /><br/>
						<label>Email:</label><br/>
						<input type="email" id="delivery_email" /><br/>
					</div>

					<div class="form-block">
						<h3>Payment</h3>
						<label>Transaction:</label><br/>
						<input type="text" id="payment_transaction" /><br/>
						<label>Request ID:</label><br/>
						<input type="text" id="payment_request_id" /><br/>
						<label>Currency:</label><br/>
						<input type="text" id="payment_currency" /><br/>
						<label>Provider:</label><br/>
						<input type="text" id="payment_provider" /><br/>
						<label>Amount:</label><br/>
						<input type="number" id="payment_amount" /><br/>
						<label>Payment DT:</label><br/>
						<input type="number" id="payment_dt" /><br/>
						<label>Bank:</label><br/>
						<input type="text" id="payment_bank" /><br/>
						<label>Delivery cost:</label><br/>
						<input type="number" id="payment_delivery_cost" /><br/>
						<label>Goods total:</label><br/>
						<input type="number" id="payment_goods_total" /><br/>
						<label>Custom fee:</label><br/>
						<input type="number" id="payment_custom_fee" /><br/>
					</div>

					<div class="form-block">
						<h3>Item</h3>
						<label>Chrt ID:</label><br/>
						<input type="number" id="item_chrt_id" /><br/>
						<label>Track number:</label><br/>
						<input type="text" id="item_track_number" /><br/>
						<label>Price:</label><br/>
						<input type="number" id="item_price" /><br/>
						<label>RID:</label><br/>
						<input type="text" id="item_rid" /><br/>
						<label>Name:</label><br/>
						<input type="text" id="item_name" /><br/>
						<label>Sale:</label><br/>
						<input type="number" id="item_sale" /><br/>
						<label>Size:</label><br/>
						<input type="text" id="item_size" /><br/>
						<label>Total price:</label><br/>
						<input type="number" id="item_total_price" /><br/>
						<label>Nm ID:</label><br/>
						<input type="number" id="item_nm_id" /><br/>
						<label>Brand:</label><br/>
						<input type="text" id="item_brand" /><br/>
						<label>Status:</label><br/>
						<input type="number" id="item_status" /><br/>
					</div>
				</div>

				<button type="submit">Add Order</button>
			</form>
			<pre id="createResult"></pre>
			<div style="flex: 1; padding: 10px; border-left: 1px solid #ccc;">
			<h2>Change order</h2>
		<form id="updateForm" onsubmit="update(); return false;">
			<div class="form-container">
				<div class="form-block">
					<h3>Orders</h3>
					<label>Order UID:</label><br/>
					<input type="text" id="oldOrderUID" required /><br/>
					<label>Track number:</label><br/>
					<input type="text" id="upd_trackNumber" required /><br/>
					<label>Entry:</label><br/>
					<input type="text" id="upd_entry" required /><br/>
					<label>Locale:</label><br/>
					<input type="text" id="upd_locale" required /><br/>
					<label>Internal signature:</label><br/>
					<input type="text" id="upd_internal_signature" required /><br/>
					<label>Customer ID:</label><br/>
					<input type="text" id="upd_customer_id" required /><br/>
					<label>Delivery service:</label><br/>
					<input type="text" id="upd_delivery_service" required /><br/>
					<label>Shardkey:</label><br/>
					<input type="text" id="upd_shardkey" required /><br/>
					<label>Sm ID:</label><br/>
					<input type="number" id="upd_sm_id" required /><br/>
					<label>Date created:</label><br/>
					<input type="text" id="upd_date_created" required /><br/>
					<label>Oof shard:</label><br/>
					<input type="text" id="upd_oof_shard" required /><br/>
				</div>
				<div class="form-block">
						<h3>Delivery</h3>
						<label>Name:</label><br/>
						<input type="text" id="delivery_name" /><br/>
						<label>Phone:</label><br/>
						<input type="text" id="delivery_phone" /><br/>
						<label>Zip:</label><br/>
						<input type="text" id="delivery_zip" /><br/>
						<label>City:</label><br/>
						<input type="text" id="delivery_city" /><br/>
						<label>Address:</label><br/>
						<input type="text" id="delivery_address" /><br/>
						<label>Region:</label><br/>
						<input type="text" id="delivery_region" /><br/>
						<label>Email:</label><br/>
						<input type="email" id="delivery_email" /><br/>
					</div>

					<div class="form-block">
						<h3>Payment</h3>
						<label>Transaction:</label><br/>
						<input type="text" id="payment_transaction" /><br/>
						<label>Request ID:</label><br/>
						<input type="text" id="payment_request_id" /><br/>
						<label>Currency:</label><br/>
						<input type="text" id="payment_currency" /><br/>
						<label>Provider:</label><br/>
						<input type="text" id="payment_provider" /><br/>
						<label>Amount:</label><br/>
						<input type="number" id="payment_amount" /><br/>
						<label>Payment DT:</label><br/>
						<input type="number" id="payment_dt" /><br/>
						<label>Bank:</label><br/>
						<input type="text" id="payment_bank" /><br/>
						<label>Delivery cost:</label><br/>
						<input type="number" id="payment_delivery_cost" /><br/>
						<label>Goods total:</label><br/>
						<input type="number" id="payment_goods_total" /><br/>
						<label>Custom fee:</label><br/>
						<input type="number" id="payment_custom_fee" /><br/>
					</div>

					<div class="form-block">
						<h3>Item</h3>
						<label>Chrt ID:</label><br/>
						<input type="number" id="item_chrt_id" /><br/>
						<label>Track number:</label><br/>
						<input type="text" id="item_track_number" /><br/>
						<label>Price:</label><br/>
						<input type="number" id="item_price" /><br/>
						<label>RID:</label><br/>
						<input type="text" id="item_rid" /><br/>
						<label>Name:</label><br/>
						<input type="text" id="item_name" /><br/>
						<label>Sale:</label><br/>
						<input type="number" id="item_sale" /><br/>
						<label>Size:</label><br/>
						<input type="text" id="item_size" /><br/>
						<label>Total price:</label><br/>
						<input type="number" id="item_total_price" /><br/>
						<label>Nm ID:</label><br/>
						<input type="number" id="item_nm_id" /><br/>
						<label>Brand:</label><br/>
						<input type="text" id="item_brand" /><br/>
						<label>Status:</label><br/>
						<input type="number" id="item_status" /><br/>
					</div>
				</div>
			</div>
			<button type="submit">Change Order</button>
		</form>
		<pre id="updateResult"></pre>

		<script>
			function getOrder() {
				let id = document.getElementById("orderId").value;
				fetch("/order/" + id)
					.then(r => {
						if (!r.ok) return r.text().then(t => { throw new Error(t) });
						return r.json();
					})
					.then(data => document.getElementById("result").textContent = JSON.stringify(data, null, 2))
					.catch(err => alert("Error: " + err));
			}

			function getAllOrders() {
				fetch("/orders")
					.then(r => {
						if (!r.ok) return r.text().then(t => { throw new Error(t) });
						return r.json();
					})
					.then(data => document.getElementById("result").textContent = JSON.stringify(data, null, 2))
					.catch(err => alert("Error: " + err));
			}

			function insert() {
				let order = {
					order_uid: document.getElementById("orderUID").value,
					track_number: document.getElementById("trackNumber").value,
					entry: document.getElementById("entry").value,
					locale: document.getElementById("locale").value,
					internal_signature: document.getElementById("internal_signature").value,
					customer_id: document.getElementById("customer_id").value,
					delivery_service: document.getElementById("delivery_service").value,
					shardkey: document.getElementById("shardkey").value,
					sm_id: parseInt(document.getElementById("sm_id").value),
					date_created: document.getElementById("date_created").value,
					oof_shard: document.getElementById("oof_shard").value,

					delivery: {
						order_uid: document.getElementById("orderUID").value,
						name: document.getElementById("delivery_name").value,
						phone: document.getElementById("delivery_phone").value,
						zip: document.getElementById("delivery_zip").value,
						city: document.getElementById("delivery_city").value,
						address: document.getElementById("delivery_address").value,
						region: document.getElementById("delivery_region").value,
						email: document.getElementById("delivery_email").value, 
					},

					payment: { 
						order_uid: document.getElementById("orderUID").value,
						transaction: document.getElementById("payment_transaction").value,
						request_id: document.getElementById("payment_request_id").value,
						currency: document.getElementById("payment_currency").value,
						provider: document.getElementById("payment_provider").value,
						amount: parseInt(document.getElementById("payment_amount").value),
						payment_dt: parseInt(document.getElementById("payment_dt").value),
						bank: document.getElementById("payment_bank").value,
						delivery_cost: parseInt(document.getElementById("payment_delivery_cost").value),
						goods_total: parseInt(document.getElementById("payment_goods_total").value),
						custom_fee: parseInt(document.getElementById("payment_custom_fee").value), 
					},
					items: [ 
						{ 
							order_uid: document.getElementById("orderUID").value,
							chrt_id: parseInt(document.getElementById("item_chrt_id").value),
							track_number: document.getElementById("item_track_number").value,
							price: parseInt(document.getElementById("item_price").value),
							rid: document.getElementById("item_rid").value,
							name: document.getElementById("item_name").value,
							sale: parseInt(document.getElementById("item_sale").value),
							size: document.getElementById("item_size").value,
							total_price: parseInt(document.getElementById("item_total_price").value),
							nm_id: parseInt(document.getElementById("item_nm_id").value),
							brand: document.getElementById("item_brand").value,
							status: parseInt(document.getElementById("item_status").value), 
						} 
					] 
					
				};
				fetch("/insert", {
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify(order)
				})
				.then(r => r.json())
				.then(data => document.getElementById("createResult").textContent = JSON.stringify(data, null, 2))
				.catch(err => alert("Error: " + err));
			}

			function update() {
				let newOrder = {
					order_uid: document.getElementById("oldOrderUID").value,
					track_number: document.getElementById("upd_trackNumber").value,
					entry: document.getElementById("upd_entry").value,
					locale: document.getElementById("upd_locale").value,
					internal_signature: document.getElementById("upd_internal_signature").value,
					customer_id: document.getElementById("upd_customer_id").value,
					delivery_service: document.getElementById("upd_delivery_service").value,
					shardkey: document.getElementById("upd_shardkey").value,
					sm_id: parseInt(document.getElementById("upd_sm_id").value),
					date_created: document.getElementById("upd_date_created").value,
					oof_shard: document.getElementById("upd_oof_shard").value,

					delivery: {
						order_uid: document.getElementById("orderUID").value,
						name: document.getElementById("delivery_name").value,
						phone: document.getElementById("delivery_phone").value,
						zip: document.getElementById("delivery_zip").value,
						city: document.getElementById("delivery_city").value,
						address: document.getElementById("delivery_address").value,
						region: document.getElementById("delivery_region").value,
						email: document.getElementById("delivery_email").value, 
					},

					payment: { 
						order_uid: document.getElementById("orderUID").value,
						transaction: document.getElementById("payment_transaction").value,
						request_id: document.getElementById("payment_request_id").value,
						currency: document.getElementById("payment_currency").value,
						provider: document.getElementById("payment_provider").value,
						amount: parseInt(document.getElementById("payment_amount").value),
						payment_dt: parseInt(document.getElementById("payment_dt").value),
						bank: document.getElementById("payment_bank").value,
						delivery_cost: parseInt(document.getElementById("payment_delivery_cost").value),
						goods_total: parseInt(document.getElementById("payment_goods_total").value),
						custom_fee: parseInt(document.getElementById("payment_custom_fee").value), 
					},
					items: [ 
						{ 
							order_uid: document.getElementById("orderUID").value,
							chrt_id: parseInt(document.getElementById("item_chrt_id").value),
							track_number: document.getElementById("item_track_number").value,
							price: parseInt(document.getElementById("item_price").value),
							rid: document.getElementById("item_rid").value,
							name: document.getElementById("item_name").value,
							sale: parseInt(document.getElementById("item_sale").value),
							size: document.getElementById("item_size").value,
							total_price: parseInt(document.getElementById("item_total_price").value),
							nm_id: parseInt(document.getElementById("item_nm_id").value),
							brand: document.getElementById("item_brand").value,
							status: parseInt(document.getElementById("item_status").value), 
						} 
					] 
				};
				fetch("/update", {
					method: "PUT",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify(newOrder),
				})
				.then(r => r.json())
				.then(data => document.getElementById("updateResult").textContent = JSON.stringify(data, null, 2))
				.catch(err => alert("Error: " + err));
			}
		</script>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
