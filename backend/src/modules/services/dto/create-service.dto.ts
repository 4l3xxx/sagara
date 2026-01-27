import { IsString, IsArray, IsOptional } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class CreateServiceDto {
  @ApiProperty()
  @IsString()
  title: string;

  @ApiProperty()
  @IsString()
  slug: string;

  @ApiProperty()
  @IsString()
  description: string;

  @ApiProperty()
  @IsString()
  shortDescription: string;

  @ApiProperty({ required: false })
  @IsOptional()
  @IsString()
  icon?: string;

  @ApiProperty({ type: [String], required: false })
  @IsOptional()
  @IsArray()
  features?: string[];

  @ApiProperty({ type: [String], required: false })
  @IsOptional()
  @IsArray()
  benefits?: string[];
}
