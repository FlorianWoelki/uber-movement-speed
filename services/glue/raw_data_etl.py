import os
import sys
from pyspark.context import SparkContext
from pyspark.sql.functions import col, upper
from awsglue.utils import getResolvedOptions
from awsglue.dynamicframe import DynamicFrame
from awsglue.context import GlueContext
from awsglue.job import Job


def main():
    sc = SparkContext.getOrCreate()
    glueContext = GlueContext(sc)
    job = Job(glueContext)

    args = getResolvedOptions(sys.argv, ["JOB_NAME"])

    job.init(args["JOB_NAME"], args)

    source_bucket = "raw-data"
    source_key = "year=2023/month=05/day=28/batch-from-d123-to-d123.csv"
    source_path = f"s3://{source_bucket}/{source_key}"
    df = glueContext.create_dynamic_frame.from_options(
        connection_type="s3",
        connection_options={"paths": [source_path]},
        format="csv",
        format_options={"withHeader": True},
    ).toDF()

    df_transformed = df.withColumn("id", col("id")).withColumn(
        "title", upper(col("title"))
    )

    df_transformed = DynamicFrame.fromDF(df_transformed, glueContext, "df_transformed")

    target_bucket = "transformed-data"
    target_key = "year=2023/month=05/day=28/"
    target_path = f"s3://{target_bucket}/{target_key}"
    glueContext.write_dynamic_frame.from_options(
        frame=df_transformed,
        connection_type="s3",
        connection_options={"path": target_path},
        format="csv",
        format_options={"quoteChar": -1},
    )

    job.commit()


main()
