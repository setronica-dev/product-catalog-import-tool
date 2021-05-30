## Mapping file
We have two mandatory fields to be filled in this file to provide ability to validate your file based on the rules.
It is ID and Category. So we need to know names of columns with values for these fields.
So it can look in the result as here:

        column-mappings:
            ID: ProductID
            Category: UNSPSC

We also use this mapping to define a mapping between field names in your source file and the connected system to transform your source file into the expected format.
