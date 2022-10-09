# Ledgerconv

Ledgerconv converts multiple bank CSV statement files into a single JSON statement file.

## How to use

Ledgerconv expects the input statement files to be arranged in the following directory structure.
```
my-input-directory
|___ name-of-your-account
|       |___ statement-file-1.csv
|       |___ statement-file-2.csv
|       |___ statement-file-3.csv
|___ for-example-hdfc-savings
|       |___ statement-file-1.csv
|       |___ statement-file-2.csv
|___ for-example-icici-credit-card
        |___ statement-file-1.csv
        |___ statement-file-2.csv
```

Once you have your statements in the correct structure, simply execute:

```shell
ledgerconv <input-dir> -o <output-dir>
```

An example can be:
```shell
ledgerconv ./statements/original -o ./statements/converted
```

The output JSON statement file contains an array of bank-agnostic transaction documents.
The schema of a transaction document is as follows:

```json
{
  "account_name": "string",
  "amount": "float64",
  "timestamp": "string",
  "bank_ref_num": "string",
  "payment_mode": "string",
  "remarks": "string"
}
```

## Supported banks

Currently, Ledgerconv can understand the CSV statement schemas of the following bank accounts:
1. ICICI savings account statements
2. ICICI credit card statements
3. HDFC savings account statements

Note that Ledgerconv uses the name of the directory (the account name) to determine the kind of parser to be used for
its statements.

See the `accountTypeInferRules` map in the `src/core/banks/mappings.go` file for more information on how account names
are parsed.

## Add support for a bank

This section describes the changes needed in the codebase to add support for a new bank account.
* Write a `ConverterFunc` for your bank account in the `src/core/banks` package. For guidance, take reference from 
  already written ConverterFunc(s).
* Go into the `src/core/banks/mappings.go` file and add your new `BankAccountType` in the list of constants.
* In the same file, update the `accountTypeInferRules` map, so that the `BankAccountType` for your new bank account can
  be inferred from an informal directory name.
* In the same file, update the `ConverterMap` map to add your new `BankAccountType` to `ConverterFunc` mapping.