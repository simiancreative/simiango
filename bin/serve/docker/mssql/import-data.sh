#run the setup script to create the DB and the schema in the DB
#do this in a loop because the timing for when the SQL instance is ready is indeterminate
/opt/mssql/bin/sqlservr --accept-eula & sleep 30;

for i in {1..50};
do
    /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P Melon40twomonkeys -d master -i setup.sql
    if [ $? -eq 0 ]
    then
        echo "setup.sql completed"
        break
    else
        echo "not ready yet..."
        sleep 1
    fi
done

#import the data from the csv file
/opt/mssql-tools/bin/bcp ExampleData.dbo.Products in \
  "/usr/src/app/Products.csv" -c -t',' -S localhost -U sa -P Melon40twomonkeys
