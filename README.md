# The tool for a product feed validation, transformation and uploading to the predefined system

## Code quality

[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=setronica-dev_product-catalog-import-tool&metric=security_rating)](https://sonarcloud.io/dashboard?id=setronica-dev_product-catalog-import-tool)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=setronica-dev_product-catalog-import-tool&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=setronica-dev_product-catalog-import-tool)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=setronica-dev_product-catalog-import-tool&metric=bugs)](https://sonarcloud.io/dashboard?id=setronica-dev_product-catalog-import-tool)

## Preamble

We are Setronica company with the great experience to build B2B and B2C integrations to exchange product catalog information between PIMs, Marketplaces, Storefronts, ERPs, services, and platforms. We know that it is quite difficult to integrate with a new service or a platform and we defined two the most important reasons:

1. Each service has its native data and API formats. It takes time to meet and configure it properly. And even if a service provides a lot of integration flexibility then it still requires to meet with this flexibility and to configure it properly;
2. Each seller has its data format with which they are used to working or just because their PIM system supports it.

This problem can be solved in three ways: on a seller side, on the system side, and by finding an external service that helps with the integration to the system. Each of them has its pros and cons. It is up to you to choose the best one for you. You may be lucky to have the integration as out of the box with the system. But all others who find it painful are welcome. This solution was built for you.

We did it in all three ways and our conclusion that the seller side is the most perspective because provides maximum flexibility and absolute concrete in what exactly should be done.
So we tried here to implement something simple to have the ability to integrate with the service quickly without diving deeply inside of it. This tool can be run locally which allows for you to get proof of the working process as soon as it is configured without any infrastructure challenges. Later it can be deployed somewhere to provide automatization and autonomy.

We implemented four steps where each of them can be switched on/off if it isn’t applicable for your case:

1. Map data from your field names to the system field names;
2. Validate data based on the system rules;
3. Transform data into the system’s native format;
4. Send data to the system.

The tool is supporting:

* the integration with the following systems: Tradeshift;
* the following formats of incoming data: CSV, XLSX;
* the following formats of outcoming data: CSV, EHF (in the nearest future).

## How does it work?

You need to do some configuration steps one time to allow it working autonomily and make automatization if you need.
After this you just need to place a file in CSV or XLSX format. Examples of files can be found in the folder 'examples'.

1. Place the source code of tool on an infrastructure and build it;
2. Place the system’s credentials into the configuration file to make a connection;
3. Configure the mapping file and. place it into the 'mapping' folder;
4. Place the file with validation rules into the 'ontology' folder;
5. Place your source file into the ‘source’ folder;
6. Run the tool [Click to see more...](./USAGE.md);

        ./product-catalog-import-tool

7. If everything is fine, then you can automate this process. Just export your files in your format into the 'source' folder and configure a scheduler to run the tool with some frequency.

## Build

You need to run the following command to build the tool and initialise default folders [Click to see more...](./INSTALL.md):

        ./install.sh -d [workdir]

where:
    [workdir] - is the path to your working directory to process files. The tool will automatically create all needed folders.

## Configuration file

This [./service.yaml](service.yaml)  file contains settings to let the tool know about all needed folders and files on the file system

        product:
          source: ./data/source/products/
          ...

This data will be filled in properly when you will build the tool. So there aren't any needs to change it manually if you don't have any special requests on it.

If you use XLSX file then the following properties has to be filled in to define a name of each sheet from your file for each type of data. It is needed in case if you use own names for these sheets. You also can define the count of header lines that has to be skipped. The configuration that you can find below means that the header line will start from third line. The first two lines will be skipped.

        xlsx_settings:
          source: ./data/source/
          sent: ./data/source/processed/
          sheet:
            products:
              name: "Products"
              header_rows_to_skip: 2
            offers:
              name: "Offers"
              header_rows_to_skip: 2
            attributes:
              name: "Attributes"
              header_rows_to_skip: 2
            offer_items:
              name: "Prices"
              header_rows_to_skip: 2


This tool contains settings to establish an API connection with the system as well.

The easiest way to get started working with the Tradeshift API is to create OAuth credentials by activating the [API Access to Own Account]( https://sandbox.tradeshift.com/#/apps/Tradeshift.AppStore/apps/Tradeshift.APIAccessToOwnAccount) app.
The app will display your credentials. Just copy these values and paste them into the configuration file:

        base_url
        consumer_key
        consumer_secret
        token
        token_secret
        tenant_id

You can also find there additional properties to define default currency and a file locale to parse correctly number delimiters in your file. Both properties are mandatory:

        currency: USD
        file_locale: en_US
    
You can also find there additional property to define a mapping between a recipient name and its tenant id on the platform. You can provide a tenant id directly in the file but it means developing some additional functionality on your export process that isn't possible for some cases or quite difficult. So you can configure it at once and apply changes when new recipients will appear to simplify it. This property is optional:

        recipients:
          - id: tenantId1
            name: tenantName1
          - id: tenantId2
            name: tenantName2

## Mapping file

This file provides ability to configure a mapping between some standard fields that known for the tool and your fields.
[Click to see more...](./MAPPING.md)

## Ontology file

This file provides ability to configure specific rules based on which all incoming files have to be validated.
[Click to see more...](./ONTOLOGY.md)
