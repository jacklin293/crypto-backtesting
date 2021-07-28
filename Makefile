dump-db-schemas:
	mysqldump -h 127.0.0.1 -u root -proot --no-data crypto_db backtests | sed -e 's/AUTO_INCREMENT=[[:digit:]]* //' > db_schemas/backtests.sql
	mysqldump -h 127.0.0.1 -u root -proot --no-data crypto_db strategies | sed -e 's/AUTO_INCREMENT=[[:digit:]]* //' > db_schemas/strategies.sql
	mysqldump -h 127.0.0.1 -u root -proot --no-data crypto_db trades | sed -e 's/AUTO_INCREMENT=[[:digit:]]* //' > db_schemas/trades.sql
	mysqldump -h 127.0.0.1 -u root -proot --no-data crypto_db moving_averages | sed -e 's/AUTO_INCREMENT=[[:digit:]]* //' > db_schemas/moving_averages.sql
