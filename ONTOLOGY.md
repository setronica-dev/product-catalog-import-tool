## Ontology file
If you have any specific rules or the system asks you to follow some rules then you can add all such rules into this file in CSV format. You will be sure that a source file that conflicts these rules never will be sent to the system.
These rules are about defining requirements for attributes. It can be specific requirements for attributes in specific categories all across whole category tree.
We are using UNSPSC classifications for now but it can be actually any that the system supports.
The file can contain the following comma separated columns:

        UNSPSC* - a category id
        UNSPSC Name - a category name
        Attribute Name* - an attribute display name
        Attribute Definition - an attribute description	
        Data Type* - data type of an attribute value. possible values: Number, Float, Text, String, Coded
        Max Character Length - max count of characters for an attribute value	
        Measurement UoM	- unit of measure of an attribute
        Is Mandatory* - an attribute is 'Mandatory' or 'Optional'
        Coded Value - a list of possible an attribute comma separated values 


        
